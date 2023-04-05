package mysql

import (
	"circle/domain"
	"gorm.io/gorm"
)

type mysqlAttendanceRepository struct {
	conn *gorm.DB
}

func NewMysqlAttendanceRepository(conn *gorm.DB) domain.AttendanceRepository {
	return &mysqlAttendanceRepository{conn}
}