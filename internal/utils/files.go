package utils

import (
	"os"
	"path/filepath"
)

func CreateTempStageFile(pattern string) (*os.File, string, error) {
	var f *os.File
	var err error
	var absPath string
	f, err = os.CreateTemp("./apsim", pattern)
	if err != nil {
		return f, absPath, err
	}
	absPath, err = filepath.Abs(f.Name())
	if err != nil {
		return f, absPath, err
	}
	return f, absPath, nil
}
