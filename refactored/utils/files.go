package utils

import (
	"apsim-api/refactored/models"
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

func GetStageFilesAbsPaths(mainCh chan models.Message) (string, string, string) {
	var csvAbsPath, constsAbsPath, apsimxAbsPath string
	for len(mainCh) > 0 {
		msg := <-mainCh
		switch msg.ID {
		case 0:
			constsAbsPath = msg.Payload
		case 1:
			csvAbsPath = msg.Payload
		case 2:
			apsimxAbsPath = msg.Payload

		}
	}
	return csvAbsPath, constsAbsPath, apsimxAbsPath
}
