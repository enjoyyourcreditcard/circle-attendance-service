package database

import (
	"circle/helper"
	"fmt"
	"net/url"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	
	// Migration
	// "circle/pkg/database/seeder"
	// "circle/domain"
)

func NewMysqlDatabase (dbUser string, dbPassword string, dbHost string, dbPort string, dbName string) *gorm.DB {
	connection 	:= fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	val 		:= url.Values{}
	val.Add("charset", "utf8mb4")
	val.Add("parseTime", "True")
	val.Add("loc", "Local")

	dsn 		:= fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	helper.PanicIfError(err)

	// Migration
	// dbConn.Migrator().DropTable(&domain.Attendance{})
	// dbConn.AutoMigrate(&domain.Attendance{})
	// seeder.Seed(dbConn)

	return dbConn
}