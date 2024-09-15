package consts

import (
	"fmt"
	"os"
	"path"
)

const (
	Available = 1
	Deleted   = 2
)

var FilePoolPath string

func PathInit() {
	RootPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	FilePoolPath = path.Join(RootPath, "filePool")
	fmt.Println("FilePoolPath automatically set to ", FilePoolPath)
	return
}

const (
	DirSize = 0
)

const (
	UserKey = "DefaultKey"
)
