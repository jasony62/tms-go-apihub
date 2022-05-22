package util

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	klog "k8s.io/klog/v2"
)

//压缩算法
type ZipCrypto struct {
	password []byte
	Keys     [3]uint32
}

func NewZipCrypto(passphrase []byte) *ZipCrypto {
	z := &ZipCrypto{}
	z.password = passphrase
	z.init()
	return z
}

func (z *ZipCrypto) init() {
	z.Keys[0] = 0x12345678
	z.Keys[1] = 0x23456789
	z.Keys[2] = 0x34567890

	for i := 0; i < len(z.password); i++ {
		z.updateKeys(z.password[i])
	}
}

func (z *ZipCrypto) updateKeys(byteValue byte) {
	z.Keys[0] = crc32update(z.Keys[0], byteValue)
	z.Keys[1] += z.Keys[0] & 0xff
	z.Keys[1] = z.Keys[1]*134775813 + 1
	z.Keys[2] = crc32update(z.Keys[2], (byte)(z.Keys[1]>>24))
}

func (z *ZipCrypto) magicByte() byte {
	var t uint32 = z.Keys[2] | 2
	return byte((t * (t ^ 1)) >> 8)
}

func (z *ZipCrypto) Encrypt(data []byte) []byte {
	length := len(data)
	chiper := make([]byte, length)
	for i := 0; i < length; i++ {
		v := data[i]
		chiper[i] = v ^ z.magicByte()
		z.updateKeys(v)
	}
	return chiper
}

func (z *ZipCrypto) Decrypt(chiper []byte) []byte {
	length := len(chiper)
	plain := make([]byte, length)
	for i, c := range chiper {
		v := c ^ z.magicByte()
		z.updateKeys(v)
		plain[i] = v
	}
	return plain
}

func crc32update(pCrc32 uint32, bval byte) uint32 {
	return crc32.IEEETable[(pCrc32^uint32(bval))&0xff] ^ (pCrc32 >> 8)
}

func ZipCryptoDecryptor(r *io.SectionReader, password []byte) (*io.SectionReader, error) {
	z := NewZipCrypto(password)
	b := make([]byte, r.Size())

	r.Read(b)

	m := z.Decrypt(b)
	return io.NewSectionReader(bytes.NewReader(m), 12, int64(len(m))), nil
}

//解压结构体
type unzip struct {
	offset int64
	fp     *os.File
	name   string
}

func (uz *unzip) init() (err error) {
	uz.fp, err = os.Open(uz.name)
	return err
}

func (uz *unzip) close() {
	if uz.fp != nil {
		uz.fp.Close()
	}
}

func (uz *unzip) Size() int64 {
	if uz.fp == nil {
		if err := uz.init(); err != nil {
			return -1
		}
	}

	fi, err := uz.fp.Stat()
	if err != nil {
		return -1
	}

	return fi.Size() - uz.offset
}

func (uz *unzip) ReadAt(p []byte, off int64) (int, error) {
	if uz.fp == nil {
		if err := uz.init(); err != nil {
			return 0, err
		}
	}

	return uz.fp.ReadAt(p, off+uz.offset)
}

func isInclude(includes []string, fname string) bool {
	if includes == nil {
		return true
	}

	for i := 0; i < len(includes); i++ {
		if includes[i] == fname {
			return true
		}
	}

	return false
}

//DeCompressZip 解压zip包
func deCompressZip(zipFile, dest, passwd string, includes []string, offset int64) error {
	uz := &unzip{offset: offset, name: zipFile}
	defer uz.close()

	klog.Infof("DeCompressZip: file:%s, dest:%s, pwd:%s\n", zipFile, dest, passwd)

	zr, err := zip.NewReader(uz, uz.Size())
	if err != nil {
		klog.Errorln("NewReader error:", err)
		return err
	}
	//如果有密码，注册接口实现带密码的压缩和解压
	if passwd != "" {
		// Register a custom Deflate compressor.
		zr.RegisterDecompressor(zip.Deflate, func(r io.Reader) io.ReadCloser {
			rs := r.(*io.SectionReader)
			r, _ = ZipCryptoDecryptor(rs, []byte(passwd))
			klog.Infoln("DeCompressZip: Register Deflate compressor")
			return flate.NewReader(r)
		})

		zr.RegisterDecompressor(zip.Store, func(r io.Reader) io.ReadCloser {
			rs := r.(*io.SectionReader)
			r, _ = ZipCryptoDecryptor(rs, []byte(passwd))
			klog.Infoln("DeCompressZip: Register Deflate decompressor")
			return ioutil.NopCloser(r)
		})
	}

	for _, f := range zr.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			klog.Errorln("MkdirAll error:", err)
			return err
		}

		inFile, err := f.Open()
		if err != nil {
			klog.Errorln("f.Open() error:", err)
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			inFile.Close()
			klog.Errorln("os.OpenFile error:", err)
			return err
		}

		_, err = io.Copy(outFile, inFile)
		inFile.Close()
		outFile.Close()
		if err != nil {
			klog.Errorln("close file error:", err)
			return err
		}
	}

	return nil
}
