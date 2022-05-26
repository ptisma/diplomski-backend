package main

import (
	"apsim-api/internal/models"
	"context"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

func main() {
	db, err := gorm.Open(sqlite.Open("baza.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	var results []models.MicroclimateReading
	counter := 0
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	ch := make(chan models.MicroclimateReading, 3)
	//blokirajuca
	//callback nije u zasebnoj gorutini
	//ako je "zapeo" u callbacku s sleepom, a context je istekao callback zavrsava i vraca se result, ali ako zapne na slanju kanalu nece izac treba select block
	result := db.Debug().WithContext(ctx).Where("location_id = ? AND date >= ? AND date <= ?", 1, "1990-01-01", "1990-03-30").FindInBatches(&results, 10, func(tx *gorm.DB, batch int) error {
		for _, result := range results {
			fmt.Println(result)
			select {
			case ch <- result:
				fmt.Println("Poslao")
			case <-ctx.Done():
				fmt.Println(ctx.Err())
				//ako ne returnam ostatak batcha ce se izvrsiti pa tek onda vratiti error u result
				return ctx.Err()
			}

		}
		counter += 1
		//if counter == 2 {
		//	time.Sleep(5 * time.Second)
		//}

		return nil
	})
	fmt.Println("Result:", result.Error)
}
