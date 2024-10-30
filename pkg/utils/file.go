package utils

import (
	"BlueNetDisk/consts"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
)

// WriteFile used when writing file to a local after the file's information is already written in to Database.
func WriteFile(file *multipart.FileHeader, dst string, uuid string) error {
	// 打开源文件
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 使用 filepath.Join 处理路径，保证跨平台兼容性
	outputPath := filepath.Join(dst, uuid+path.Ext(file.Filename))

	// 创建目标文件
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 复制文件内容
	_, err = io.Copy(out, src)
	if err != nil {
		// 如果复制失败，删除目标文件
		_ = os.Remove(outputPath)
		return err
	}
	return nil
}

func RemoveFile(filename string) error {
	return os.Remove(path.Join(consts.FilePoolPath, filename))
}

func Sha1(file *multipart.FileHeader) string {
	hash := sha1.New()
	fd, err := file.Open()
	if err != nil {
		return ""
	}
	defer fd.Close()
	_, err = io.Copy(hash, fd)
	if err != nil {
		return ""
	}
	v := hex.EncodeToString(hash.Sum(nil))
	return v
}
