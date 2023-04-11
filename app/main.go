package main

import (
	"circle/helper"
	"circle/pkg/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	_attendanceHttpDelivery "circle/infrastructure/attendance/delivery/http"
	_attendanceUsecase "circle/infrastructure/attendance/usecase"
	_attendanceRepo "circle/infrastructure/attendance/repository/mysql"
)

func init()  {
		   viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	helper.PanicIfError(err)
}

func main() {
	dbUser 			:= viper.GetString(`database.user`)
	dbPassword 		:= viper.GetString(`database.pass`)
	dbHost 			:= viper.GetString(`database.host`)
	dbPort 			:= viper.GetString(`database.port`)
	dbName 			:= viper.GetString(`database.name`)
	portService 	:= viper.GetString(`server.address`)
	dbConn 			:= database.NewMysqlDatabase(dbUser, dbPassword, dbHost, dbPort, dbName)
	timeoutContext 	:= time.Duration(viper.GetInt("context.timeout")) * time.Second

	app := fiber.New(fiber.Config{
			Prefork:       	true,
			CaseSensitive: 	true,
			StrictRouting: 	true,
			ServerHeader:  	"Fiber",
			AppName: 		"Circle",
	})

	// Init Repository
	atr := _attendanceRepo.NewMysqlAttendanceRepository(dbConn)

	// Init Usecase
	atu := _attendanceUsecase.NewAttendanceUsecase(atr, timeoutContext)

	// Init Delivery
	_attendanceHttpDelivery.NewAttendanceHandler(app, atu)
	
	err := app.Listen(portService)
	helper.PanicIfError(err)
}