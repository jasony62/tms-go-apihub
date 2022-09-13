package logger

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"net"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const DefaultLogPath = "/var/log/test" // 默认输出日志文件路径

type LogConfigs struct {
	LogLevel          string // 日志打印级别 debug  info  warning  error
	LogFormat         string // 输出日志格式	logfmt, json
	LogPath           string // 输出日志文件路径
	LogFileName       string // 输出日志文件名称
	LogFileMaxSize    int    // 【日志分割】单个日志文件最多存储量 单位(mb)
	LogFileMaxBackups int    // 【日志分割】日志备份文件最多数量
	LogMaxAge         int    // 日志保留时间，单位: 天 (day)
	LogCompress       bool   // 是否压缩日志
	LogStdout         bool   // 是否输出到控制台
}

// InitLogger 初始化 log
func InitLogger(conf LogConfigs, logwithlevel bool) error {
	encoder := getEncoder(conf) // 获取日志输出编码

	var coreArr []zapcore.Core
	var level zapcore.Level
	if conf.LogLevel == "info" {
		level = zap.InfoLevel
	} else if conf.LogLevel == "warn" {
		level = zap.WarnLevel
	} else if conf.LogLevel == "error" {
		level = zap.ErrorLevel
	} else {
		level = zap.DebugLevel
	}
	if logwithlevel {
		errorLog := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别及以上
			return lev >= zap.ErrorLevel
		})
		allLog := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //所有级别同时输出
			return lev <= zap.FatalLevel && lev >= level
		})

		allFileWriteSyncer, err := getLogWriter("all", conf) // 日志文件配置 文件位置和切割
		if err != nil {
			LogS().Errorln("设定日志info输出方式失败:", err.Error())
			return err
		}

		errorFileWriteSyncer, err := getLogWriter("error", conf) // 日志文件配置 文件位置和切割
		if err != nil {
			LogS().Errorln("设定日志error输出方式失败:", err.Error())
			return err
		}
		//第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
		allFileCore := zapcore.NewCore(encoder, allFileWriteSyncer, allLog)
		errorFileCore := zapcore.NewCore(encoder, errorFileWriteSyncer, errorLog)

		coreArr = append(coreArr, allFileCore)
		coreArr = append(coreArr, errorFileCore)
	} else {
		allLog := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //所有级别日志
			return lev <= zap.FatalLevel && lev >= level
		})

		allFileWriteSyncer, err := getLogWriter("", conf) // 日志文件配置 文件位置和切割
		if err != nil {
			LogS().Errorln("设定日志输出方式失败:", err.Error())
			return err
		}

		allFileCore := zapcore.NewCore(encoder, allFileWriteSyncer, allLog)

		coreArr = append(coreArr, allFileCore)
	}

	logging := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller()) //zap.AddCaller()为显示文件名和行号，可省略

	zap.ReplaceGlobals(logging)
	return nil
}

// getEncoder 编码器(如何写入日志)
func getEncoder(conf LogConfigs) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // log 时间格式 例如: 2021-09-11t20:05:54.852+0800
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if conf.LogFormat == "json" {
		return zapcore.NewJSONEncoder(encoderConfig) // 以json格式写入
	}
	return zapcore.NewConsoleEncoder(encoderConfig) // 以logfmt格式写入
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// getLogWriter 获取日志输出方式  日志文件 控制台
func getLogWriter(level string, conf LogConfigs) (zapcore.WriteSyncer, error) {
	// 判断日志路径是否存在，如果不存在就创建
	if exist, _ := pathExists(conf.LogPath); !exist {
		if conf.LogPath == "" {
			conf.LogPath = DefaultLogPath
		}
		if err := os.MkdirAll(conf.LogPath, os.ModePerm); err != nil {
			conf.LogPath = DefaultLogPath
			if err := os.MkdirAll(conf.LogPath, os.ModePerm); err != nil {
				return nil, err
			}
		}
	}

	//判断文件是否存在，存在则删除
	var logfile string
	if level != "" {
		logfile = conf.LogPath + level + "/" + conf.LogFileName
	} else {
		logfile = conf.LogPath + conf.LogFileName
	}

	if exist, _ := pathExists(logfile); exist {
		os.Remove(logfile)
	}

	// 日志文件 与 日志切割 配置
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(conf.LogPath, level, conf.LogFileName), // 日志文件路径
		MaxSize:    conf.LogFileMaxSize,                                  // 单个日志文件最大多少 mb
		MaxBackups: conf.LogFileMaxBackups,                               // 日志备份数量
		MaxAge:     conf.LogMaxAge,                                       // 日志最长保留时间
		Compress:   conf.LogCompress,                                     // 是否压缩日志
	}

	//日志输出到控制台
	if len(conf.LogFileName) == 0 {
		return zapcore.AddSync(os.Stdout), nil
	}
	if conf.LogStdout {
		// 日志同时输出到控制台和日志文件中
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout)), nil
	} else {
		// 日志只输出到日志文件
		return zapcore.AddSync(lumberJackLogger), nil
	}
}

func LogS() *zap.SugaredLogger {
	return zap.S()
}

func LogL() *zap.Logger {
	return zap.L()
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
