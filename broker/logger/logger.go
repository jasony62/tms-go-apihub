package logger

import (
	"os"
	"path/filepath"

	"github.com/jasony62/tms-go-apihub/tool"
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
	if logwithlevel {
		//日志级别
		highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别及以上
			return lev >= zap.ErrorLevel
		})
		lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //所有级别同时输出
			return lev <= zap.FatalLevel && lev >= zap.DebugLevel
		})

		infoFileWriteSyncer, err := getLogWriter("all", conf) // 日志文件配置 文件位置和切割
		if err != nil {
			Errorln("设定日志info输出方式失败:", err.Error())
			return err
		}

		errorFileWriteSyncer, err := getLogWriter("error", conf) // 日志文件配置 文件位置和切割
		if err != nil {
			Errorln("设定日志error输出方式失败:", err.Error())
			return err
		}
		//第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
		infoFileCore := zapcore.NewCore(encoder, infoFileWriteSyncer, lowPriority)

		//第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
		errorFileCore := zapcore.NewCore(encoder, errorFileWriteSyncer, highPriority)

		coreArr = append(coreArr, infoFileCore)
		coreArr = append(coreArr, errorFileCore)
	} else {
		lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
			return lev <= zap.FatalLevel && lev >= zap.DebugLevel
		})

		infoFileWriteSyncer, err := getLogWriter("", conf) // 日志文件配置 文件位置和切割
		if err != nil {
			Errorln("设定日志输出方式失败:", err.Error())
			return err
		}

		//第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
		infoFileCore := zapcore.NewCore(encoder, infoFileWriteSyncer, lowPriority)

		coreArr = append(coreArr, infoFileCore)
	}

	logger := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller()) //zap.AddCaller()为显示文件名和行号，可省略

	zap.ReplaceGlobals(logger)

	return nil
}

// getEncoder 编码器(如何写入日志)
func getEncoder(conf LogConfigs) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // log 时间格式 例如: 2021-09-11t20:05:54.852+0800
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 输出level序列化为全大写字符串，如 INFO DEBUG ERROR
	//encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	//encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if conf.LogFormat == "json" {
		return zapcore.NewJSONEncoder(encoderConfig) // 以json格式写入
	}
	return zapcore.NewConsoleEncoder(encoderConfig) // 以logfmt格式写入
}

// getLogWriter 获取日志输出方式  日志文件 控制台
func getLogWriter(level string, conf LogConfigs) (zapcore.WriteSyncer, error) {

	// 判断日志路径是否存在，如果不存在就创建
	if exist, _ := tool.PathExists(conf.LogPath); !exist {
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

	if exist, _ := tool.PathExists(logfile); exist {
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

func Logger() *zap.SugaredLogger {
	return zap.S()
}

func Debugln(template string, args ...interface{}) {
	zap.S().Debugln(template, args)
}

func Infoln(template string, args ...interface{}) {
	zap.S().Infoln(template, args)
}

func Warningln(template string, args ...interface{}) {
	zap.S().Warnln(template, args)
}

func Errorln(template string, args ...interface{}) {
	zap.S().Errorln(template, args)
}

func DPanicln(template string, args ...interface{}) {
	zap.S().DPanicln(template, args)
}

func Panicln(template string, args ...interface{}) {
	zap.S().Panicln(template, args)
}

func Fatalln(template string, args ...interface{}) {
	zap.S().Fatalln(template, args)
}
