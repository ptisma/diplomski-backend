package utils

import (
	"os"
	"path/filepath"
)

// Create temp file (random name) inside apsim-stage-area dir
func CreateTempStageFile(pattern string) (*os.File, string, error) {
	var f *os.File
	var err error
	var absPath string
	f, err = os.CreateTemp("./apsim-stage-area", pattern)
	if err != nil {
		return nil, absPath, err
	}
	absPath, err = filepath.Abs(f.Name())
	if err != nil {
		return nil, absPath, err
	}
	return f, absPath, err
}

// Receive absolute paths of the newly created files using chan (goroutines send their abs path to the said chan)
func GetStageFilesAbsPaths(mainCh chan Message) (string, string, string) {
	defer close(mainCh)
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
