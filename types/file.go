package types

// UploadFileResp is used when fail to up load a file,info saves error messages
type UploadFileResp struct {
	FileName string
	FileSize int64
	Info     string
}
