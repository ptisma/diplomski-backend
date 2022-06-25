package utils

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
)

// Run ApsimX simulation using dll inside apsim-cli dir
func RunAPSIMSimulation(ctx context.Context, absPathApsimx string) error {

	cmd := exec.CommandContext(ctx, "dotnet", "apsim.dll", "run", "--single-threaded", "f", absPathApsimx)
	//sad sam u apsimu
	cmd.Dir = "./apsim-cli"
	err := cmd.Run()
	return err
}

// Construct a absolute path of the newly created .db file after the simulation
func ConstructDBAbsPath(absPathApsimx string) string {
	baseF := filepath.Base(absPathApsimx)
	fNoExt := baseF[:len(baseF)-len(filepath.Ext(baseF))]
	dbFile := filepath.Join(filepath.Dir(absPathApsimx), fNoExt+".db")

	return dbFile
}

// Delete a file based on its path
func DeleteStageFile(absPathFile string) error {

	return os.Remove(absPathFile)
}
