package main

import (
	"circle/helper"
	"circle/pkg/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	_attendanceHttpDelivery "circle/infrastructure/attendance/delivery/http"
	_apimysqlAttendanceRepo "circle/infrastructure/attendance/repository/api"
	_mysqlmysqlAttendanceRepo "circle/infrastructure/attendance/repository/mysql"
	_attendanceUsecase "circle/infrastructure/attendance/usecase"

	_assignmentHttpDelivery "circle/infrastructure/assignment/delivery/http"
	_assignmentRepo "circle/infrastructure/assignment/repository/mysql"
	_assignmentUsecase "circle/infrastructure/assignment/usecase"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	helper.PanicIfError(err)
}

func main() {
	dbUser := viper.GetString(`database.user`)
	dbPassword := viper.GetString(`database.pass`)
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbName := viper.GetString(`database.name`)
	portService := viper.GetString(`server.address`)
	dbConn := database.NewMysqlDatabase(dbUser, dbPassword, dbHost, dbPort, dbName)
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Circle",
		BodyLimit:     2097152,
	})

	var baseUrl string
	if viper.GetString(`server.name`) == "local" {
		baseUrl = "http://localhost"
	} else {
		baseUrl = "https://istudio.mncplay.id"

	}

	// Init Repository
	atr := _mysqlmysqlAttendanceRepo.NewMysqlmysqlAttendanceRepository(dbConn)
	ity := _apimysqlAttendanceRepo.NewApimysqlAttendanceRepository(baseUrl)

	asr := _assignmentRepo.NewMysqlAssignmentRepository(dbConn)

	// Init Usecase
	atu := _attendanceUsecase.NewAttendanceUsecase(atr, ity, timeoutContext)
	asu := _assignmentUsecase.NewAssignmentUsecase(asr, timeoutContext)

	// Init Delivery
	_attendanceHttpDelivery.NewAttendanceHandler(app, atu)
	_assignmentHttpDelivery.NewAssignmentHandler(app, asu)

	err := app.Listen(portService)
	helper.PanicIfError(err)
}
