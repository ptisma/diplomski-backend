package utils

import (
	"apsim-api/internal/models"
	"apsim-api/pkg/application"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"os/exec"
	"path/filepath"
)

func RunAPSIMSimulation(absPathApsimx string) error {

	cmd := exec.Command("dotnet", "apsim.dll", "run", "--single-threaded", "f", absPathApsimx)
	//sad sam u apsimu
	cmd.Dir = "../apsim-cli/netcoreapp3.1"
	err := cmd.Run()

	return err
}

func ReadAPSIMSimulationResults(dbFilePath string, app *application.Application) ([]models.Yield, error) {

	var yields = []models.Yield{}
	fmt.Println("Opening db file:", dbFilePath)
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		return yields, err
	}
	err = db.Raw(`SELECT strftime('%Y', date) as year, max(yield) as yield FROM report GROUP BY year`).Scan(&yields).Error
	fmt.Println("yields:", yields)

	return yields, err

}

func ConstructDBAbsPath(absPathApsimx string) string {
	baseF := filepath.Base(absPathApsimx)
	fNoExt := baseF[:len(baseF)-len(filepath.Ext(baseF))]
	dbFile := filepath.Join(filepath.Dir(absPathApsimx), fNoExt+".db")

	return dbFile
}

func DeleteStageFile(absPathFile string) error {

	return os.Remove(absPathFile)
}
