package mysql

import (
	"circle/domain"
	"circle/infrastructure/attendance/helper"
	"context"

	"gorm.io/gorm"
)

type mysqlAttendanceRepository struct {
	conn *gorm.DB
}

func NewMysqlAttendanceRepository(conn *gorm.DB) domain.AttendanceRepository {
	return &mysqlAttendanceRepository{conn}
}

func (ar mysqlAttendanceRepository) GetUserLastAttendance(ctx context.Context, userId string) (domain.Attendance, error) {
	var attendance domain.Attendance

	result := ar.conn.Where("user_id = ?", userId).Last(&attendance)
	return attendance, result.Error
}

func (ar mysqlAttendanceRepository) GetUserAttendanceMonthly(ctx context.Context, formatedCurrentMY string, userId string) (domain.AttendanceMonthly, error) {
	var attendances 		[]domain.Attendance
	var attendanceMonthly 	domain.AttendanceMonthly
	var totalAbsen 			int64
	var totalWfo 			int64
	var totalWfh 			int64
	var totalOff 			int64

	result := ar.conn.Where("user_id = ?", userId).Where("start_at LIKE ?", "%"+formatedCurrentMY+"%").Find(&attendances).Count(&totalAbsen)
	ar.conn.Where("user_id = ?", userId).Where("start_at LIKE ?", "%"+formatedCurrentMY+"%").Where("worktype = ?", "wfo").Find(&attendances).Count(&totalWfo)
	ar.conn.Where("user_id = ?", userId).Where("start_at LIKE ?", "%"+formatedCurrentMY+"%").Where("worktype = ?", "wfh").Find(&attendances).Count(&totalWfh)
	ar.conn.Where("user_id = ?", userId).Where("start_at LIKE ?", "%"+formatedCurrentMY+"%").Where("worktype = ?", "off").Find(&attendances).Count(&totalOff)

	attendanceMonthly = helper.GetUserAttendanceMonthly(attendanceMonthly, userId, totalAbsen, totalWfo, totalWfh, totalOff)
	return attendanceMonthly, result.Error
}

func (ar mysqlAttendanceRepository) CheckAbsen(ctx context.Context, userId string, formatedCurrentDate string) (int, error) {
	var counted int64
	var Attendance domain.Attendance

	err := ar.conn.Where("user_id = ?", userId).Where("start_at LIKE ?", formatedCurrentDate+"%").Find(&Attendance).Count(&counted)
	result := int(counted)
	return result, err.Error
}

func (ar mysqlAttendanceRepository) CreateAbsen(ctx context.Context, attendance *domain.Attendance) (*domain.Attendance, error) {
	result := ar.conn.Create(attendance)
	return attendance, result.Error
}

func (ar mysqlAttendanceRepository) UpdateAbsen(ctx context.Context, endAttendance *domain.EndAttendance, attendanceId int) error {
	var attendance domain.Attendance

	result := ar.conn.Model(attendance).Where("id = ?", attendanceId).Updates(*endAttendance)
	return result.Error
}

func (ar mysqlAttendanceRepository) Hello() string {
	return "hello"
}
