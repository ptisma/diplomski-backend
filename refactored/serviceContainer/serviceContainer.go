package serviceContainer

import (
	"apsim-api/internal/models"
	"apsim-api/refactored/controllers"
	"apsim-api/refactored/repositories"
	"apsim-api/refactored/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
)

type IServiceContainer interface {
	InjectCultureController() controllers.CultureController
	InjectLocationController() controllers.LocationController
	InjectMicroclimateController() controllers.MicroclimateController
	InjectMicroclimateReadingController() controllers.MicroclimateReadingController
	InjectYieldController() controllers.YieldController
}

type kernel struct{}

func (k *kernel) InjectCultureController() controllers.CultureController {

	db, _ := gorm.Open(sqlite.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-api\\baza.db"), &gorm.Config{})
	db.AutoMigrate(models.Culture{})

	//if err != nil {
	//	return controllers.CultureController{}, err
	//}

	cultureRepository := &repositories.CultureRepository{DB: db}

	cultureService := &services.CultureService{cultureRepository}

	cultureController := controllers.CultureController{cultureService}

	return cultureController
}

func (k *kernel) InjectLocationController() controllers.LocationController {

	db, _ := gorm.Open(sqlite.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-api\\baza.db"), &gorm.Config{})
	db.AutoMigrate(models.Location{})

	//if err != nil {
	//	return controllers.CultureController{}, err
	//}

	locationRepository := &repositories.LocationRepository{DB: db}

	locationService := &services.LocationService{locationRepository}

	locationController := controllers.LocationController{locationService}

	return locationController
}

func (k *kernel) InjectMicroclimateController() controllers.MicroclimateController {

	db, _ := gorm.Open(sqlite.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-api\\baza.db"), &gorm.Config{})
	db.AutoMigrate(models.Microclimate{})

	//if err != nil {
	//	return controllers.CultureController{}, err
	//}

	microclimateRepository := &repositories.MicroclimateRepository{DB: db}

	microclimateService := &services.MicroclimateService{microclimateRepository}

	microclimateController := controllers.MicroclimateController{microclimateService}

	return microclimateController
}

func (k *kernel) InjectMicroclimateReadingController() controllers.MicroclimateReadingController {

	db, _ := gorm.Open(sqlite.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-api\\baza.db"), &gorm.Config{})
	db.AutoMigrate(models.MicroclimateReading{})

	//if err != nil {
	//	return controllers.CultureController{}, err
	//}

	microclimateReadingRepository := &repositories.MicroclimateReadingRepository{DB: db}
	predictedMicroclimateReadingRepository := &repositories.PredictedMicroclimateReadingRepository{DB: db}

	microclimateReadingService := &services.MicroclimateReadingService{microclimateReadingRepository, predictedMicroclimateReadingRepository}

	microclimateReadingController := controllers.MicroclimateReadingController{microclimateReadingService}

	return microclimateReadingController
}

func (k *kernel) InjectYieldController() controllers.YieldController {

	db, _ := gorm.Open(sqlite.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-api\\baza.db"), &gorm.Config{})

	soilRepository := &repositories.SoilRepository{DB: db}

	soilService := &services.SoilService{soilRepository}

	cultureRepository := &repositories.CultureRepository{DB: db}

	cultureService := &services.CultureService{cultureRepository}

	locationRepository := &repositories.LocationRepository{DB: db}

	locationService := &services.LocationService{locationRepository}

	microclimateReadingRepository := &repositories.MicroclimateReadingRepository{DB: db}
	predictedMicroclimateReadingRepository := &repositories.PredictedMicroclimateReadingRepository{DB: db}

	microclimateReadingService := &services.MicroclimateReadingService{microclimateReadingRepository, predictedMicroclimateReadingRepository}

	//if err != nil {
	//	return controllers.CultureController{}, err
	//}

	yieldController := controllers.YieldController{
		LocationService:            locationService,
		CultureService:             cultureService,
		MicroclimateReadingService: microclimateReadingService,
		SoilService:                soilService,
	}

	return yieldController

}

var (
	k             *kernel
	containerOnce sync.Once
)

func ServiceContainer() IServiceContainer {
	if k == nil {
		containerOnce.Do(func() {
			k = &kernel{}
		})
	}
	return k
}
