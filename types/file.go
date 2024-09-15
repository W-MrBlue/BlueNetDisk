package types

import (
	"time"
)

// UploadFileResp is used when fail to upload a file,info saves error messages
type UploadFileResp struct {
	FileName string `json:"filename"`
	FileSize int64  `json:"filesize"`
	Info     string `json:"info"`
}

type ListFileReq struct {
	ParentId string `json:"parent_id"`
}
type ListFileResp struct {
	UUID      string `json:"uuid"`
	FileName  string `json:"filename"`
	FileSize  int64  `json:"filesize"`
	Ext       string `json:"ext"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type MkdirReq struct {
	DirName  string `json:"dirname"`
	ParentId string `json:"parent_id"`
}

type DeleteFileReq struct {
	FileName string `json:"filename"`
	ParentId string `json:"parent_id"`
}

type RenameFileReq struct {
	NewFileName string `json:"new_filename"`
	FileId      string `json:"file_id"`
	ParentId    string `json:"parent_id"`
}

type DownloadFileReq struct {
	FileName string `json:"filename"`
	ParentId string `json:"parent_id"`
}
