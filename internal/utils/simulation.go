package utils

import (
	"apsim-api/internal/models"
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
	cmd.Dir = "./apsim-cli"
	err := cmd.Run()

	return err
}

func ReadAPSIMSimulationResults(dbFilePath string) ([]models.Yield, error) {

	var yields = []models.Yield{}
	fmt.Println("Opening db file:", dbFilePath)
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		return yields, err
	}
	err = db.Raw(`SELECT strftime('%Y', date) as year, max(yield) as yield FROM report GROUP BY year`).Scan(&yields).Error

	//for _, j := range yields {
	//	err = db.Raw(`SELECT date FROM report`).Scan(&yields).Error
	//	j.Dates =
	//
	//}
	//fmt.Println("yields:", yields)

	client, _ := db.DB()
	//not handled
	client.Close()

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
