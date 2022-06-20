package serviceContainer

import (
	"apsim-api/internal/controllers"
	"apsim-api/internal/infra/application"
	"apsim-api/internal/repositories"
	"apsim-api/internal/services"
	"sync"
)

type IServiceContainer interface {
	InjectCultureController() controllers.CultureController
	InjectLocationController() controllers.LocationController
	InjectMicroclimateController() controllers.MicroclimateController
	InjectMicroclimateReadingController() controllers.MicroclimateReadingController
	InjectYieldController() controllers.YieldController
	InjectGrowingDegreeDayController() controllers.GrowingDegreeDayController
}

type skernel struct {
	app application.Application
}

func (k *skernel) InjectCultureController() controllers.CultureController {

	//db, _ := gorm.Open(sqlite.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-stage-area-api\\baza.db"), &gorm.Config{})
	//db.AutoMigrate(models.Culture{})

	//if err != nil {
	//	return controllers.CultureController{}, err
	//}

	cultureRepository := &repositories.CultureRepository{DB: k.app.GetDB().GetClient()}

	cultureService := &services.CultureService{cultureRepository}

	cultureController := controllers.CultureController{cultureService}

	return cultureController
}

func (k *skernel) InjectLocationController() controllers.LocationController {

	//db, _ := gorm.Open(sqlite.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-stage-area-api\\baza.db"), &gorm.Config{})
	//db.AutoMigrate(models.Location{})

	//if err != nil {
	//	return controllers.CultureController{}, err
	//}

	locationRepository := &repositories.LocationRepository{DB: k.app.GetDB().GetClient()}

	locationService := &services.LocationService{locationRepository}

	locationController := controllers.LocationController{locationService}

	return locationController
}

func (k *skernel) InjectMicroclimateController() controllers.MicroclimateController {

	//db, _ := gorm.Open(sqlite.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-stage-area-api\\baza.db"), &gorm.Config{})
	//db.AutoMigrate(models.Microclimate{})

	//if err != nil {
	//	return controllers.CultureController{}, err
	//}

	microclimateRepository := &repositories.MicroclimateRepository{DB: k.app.GetDB().GetClient()}

	microclimateService := &services.MicroclimateService{microclimateRepository}

	microclimateController := controllers.MicroclimateController{microclimateService}

	return microclimateController
}

func (k *skernel) InjectMicroclimateReadingController() controllers.MicroclimateReadingController {

	//db, _ := gorm.Open(sqlite.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-stage-area-api\\baza.db"), &gorm.Config{})
	//db.AutoMigrate(models.MicroclimateReading{})

	//if err != nil {
	//	return controllers.CultureController{}, err
	//}

	microclimateReadingRepository := &repositories.MicroclimateReadingRepository{DB: k.app.GetDB().GetClient()}
	predictedMicroclimateReadingRepository := &repositories.PredictedMicroclimateReadingRepository{DB: k.app.GetDB().GetClient()}

	microclimateReadingService := &services.MicroclimateReadingService{microclimateReadingRepository, predictedMicroclimateReadingRepository}

	microclimateReadingController := controllers.MicroclimateReadingController{microclimateReadingService}

	return microclimateReadingController
}

func (k *skernel) InjectYieldController() controllers.YieldController {

	//db, _ := gorm.Open(sqlite.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-stage-area-api\\baza.db"), &gorm.Config{})
	//
	//client := influxdb2.NewClient("http://localhost:8086", "iCIavmWDn08-O9C5R4qjR2xrWN-57YluKNJY6HW6NCBEKPXMy_AdwwFmIi0k5TDWKRdkT6f2P4Wpe4QhVOJExQ==")

	yieldRepository := &repositories.YieldRepository{
		Client:      k.app.GetCache().GetClient(),
		Org:         k.app.GetConfig().GetInfluxDbOrg(),
		Bucket:      k.app.GetConfig().GetInfluxDbBucket(),
		Measurement: k.app.GetConfig().GetInfluxDbMeasurement(),
	}

	yieldService := &services.YieldService{yieldRepository}

	soilRepository := &repositories.SoilRepository{DB: k.app.GetDB().GetClient()}

	soilService := &services.SoilService{soilRepository}

	cultureRepository := &repositories.CultureRepository{DB: k.app.GetDB().GetClient()}

	cultureService := &services.CultureService{cultureRepository}

	locationRepository := &repositories.LocationRepository{DB: k.app.GetDB().GetClient()}

	locationService := &services.LocationService{locationRepository}

	microclimateReadingRepository := &repositories.MicroclimateReadingRepository{DB: k.app.GetDB().GetClient()}
	predictedMicroclimateReadingRepository := &repositories.PredictedMicroclimateReadingRepository{DB: k.app.GetDB().GetClient()}

	microclimateReadingService := &services.MicroclimateReadingService{microclimateReadingRepository, predictedMicroclimateReadingRepository}

	//if err != nil {
	//	return controllers.CultureController{}, err
	//}

	yieldController := controllers.YieldController{
		LocationService:            locationService,
		CultureService:             cultureService,
		MicroclimateReadingService: microclimateReadingService,
		SoilService:                soilService,
		YieldService:               yieldService,
	}

	return yieldController

}

func (k *skernel) InjectGrowingDegreeDayController() controllers.GrowingDegreeDayController {
	//db, _ := gorm.Open(sqlite.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-stage-area-api\\baza.db"), &gorm.Config{})
	//db.AutoMigrate(models.Location{})

	//if err != nil {
	//	return controllers.CultureController{}, err
	//}
	cultureRepository := &repositories.CultureRepository{DB: k.app.GetDB().GetClient()}

	cultureService := &services.CultureService{cultureRepository}

	microclimateReadingRepository := &repositories.MicroclimateReadingRepository{DB: k.app.GetDB().GetClient()}
	predictedMicroclimateReadingRepository := &repositories.PredictedMicroclimateReadingRepository{DB: k.app.GetDB().GetClient()}

	microclimateReadingService := &services.MicroclimateReadingService{microclimateReadingRepository, predictedMicroclimateReadingRepository}
	microclimateRepository := &repositories.MicroclimateRepository{DB: k.app.GetDB().GetClient()}

	microclimateService := &services.MicroclimateService{microclimateRepository}

	growingDegreeDayController := controllers.GrowingDegreeDayController{cultureService, microclimateReadingService, microclimateService}

	return growingDegreeDayController
}

var (
	k             *skernel
	containerOnce sync.Once
)

func GetServiceContainer(app application.Application) IServiceContainer {
	if k == nil {
		containerOnce.Do(func() {
			k = &skernel{app: app}
		})
	}
	return k

}
