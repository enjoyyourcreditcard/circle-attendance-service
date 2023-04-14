package mysql

import (
	"circle/domain"
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

func (ar mysqlAttendanceRepository) GetUserAttendanceData(ctx context.Context, userId string, startAt string, endAt string) ([]domain.Attendance, error) {
	var attendances []domain.Attendance
	result := ar.conn.Where("user_id = ?", userId).Where("start_at BETWEEN ? AND ?", startAt, endAt).Where("end_at BETWEEN ? AND ?", startAt, endAt).Find(&attendances)
	return attendances, result.Error
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

func (ar mysqlAttendanceRepository) PostAttendanceNotes(ctx context.Context, userId string, notes string) error {
	var attendance domain.Attendance

	result := ar.conn.Where("user_id = ?", userId).Last(&attendance).Update("notes", notes)
	return result.Error
}

func (ar mysqlAttendanceRepository) CreateAssignment(ctx context.Context, assignment *domain.Assignment) error {
	return nil
}
