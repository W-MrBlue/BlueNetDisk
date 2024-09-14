package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path"
)

// Writer used when writing file to a local after the file's information is already written in to Database.
// db(*gorm.DB) is passed as params to rollback db operations
func Writer(file *multipart.FileHeader, dst string, uuid string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst + "\\" + uuid + path.Ext(file.Filename))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		_ = os.Remove(dst)
		return err
	}
	return nil
}

func Sha1(file *multipart.FileHeader) string {
	hash := sha1.New()
	fd, err := file.Open()
	if err != nil {
		return ""
	}
	defer fd.Close()
	data := make([]byte, 0)
	_, err = fd.Read(data)
	return hex.EncodeToString(hash.Sum(data))
}
