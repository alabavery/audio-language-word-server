package util

import (
	"encoding/json"
	"os"
)

func ReadFileThatMayNotExist(filePath string) ([]byte, bool) {
	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		return []byte{}, false
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	bytes := make([]byte, fileInfo.Size())
	file.Read(bytes)
	return bytes, true
}

// GetJSONWhenFileMayNotExist saves content to saveTo. Returns whether or not
// file existed.  saves nothing to saveTo if file did not exist
func GetJSONWhenFileMayNotExist(filePath string, saveTo interface{}) bool {
	bytes, exists := ReadFileThatMayNotExist(filePath)
	if exists == false {
		return false
	}
	err := json.Unmarshal(bytes, &saveTo)
	if err != nil {
		panic(err)
	}
	return true
}
