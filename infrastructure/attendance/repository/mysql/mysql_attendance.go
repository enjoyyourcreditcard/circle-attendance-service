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

func (ar mysqlAttendanceRepository) GetUserDashboardAttendance(ctx context.Context, userId string, startAt string, endAt string) (domain.DashboardAttendance, error) {
	var attendance 			domain.Attendance
	var dashboardAttendance domain.DashboardAttendance
	
	result := ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Count(&dashboardAttendance.WorkingDay)
	result = ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Count(&dashboardAttendance.TotalClockin)
	result = ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Not("end_at LIKE ?", "%00:00:00").Count(&dashboardAttendance.TotalClockout)
	result = ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Where("worktype = ?", "wfh").Count(&dashboardAttendance.TotalWfh)
	result = ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Where("worktype = ?", "wfo").Count(&dashboardAttendance.TotalWfo)
	result = ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Where("TIME(STR_TO_DATE(start_at, '%d-%m-%Y %H:%i:%s')) > ?", "09:15:00").Count(&dashboardAttendance.LateIn)
	result = ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Where("TIME(STR_TO_DATE(start_at, '%d-%m-%Y %H:%i:%s')) < ?", "09:15:00").Count(&dashboardAttendance.EarlyIn)
	result = ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Where("TIME(STR_TO_DATE(end_at, '%d-%m-%Y %H:%i:%s')) < ?", "18:00:00").Not("TIME(STR_TO_DATE(end_at, '%d-%m-%Y %H:%i:%s')) = ?", "00:00:00").Count(&dashboardAttendance.EarlyOut)
	result = ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Where("status_start = ?", "inside_area").Count(&dashboardAttendance.InsideArea)
	result = ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Not("status_start != ?", "inside_area").Count(&dashboardAttendance.OutsideArea)
	result = ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Where("status_start = ?", "inside_other_area").Count(&dashboardAttendance.InsideOtherArea)
	result = ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Where("type = ?", "shifting").Count(&dashboardAttendance.Shifting)
	result = ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Not("type = ?", "shifting").Count(&dashboardAttendance.OfficeHour)
	result = ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Not("working_hour = ?", "").Where("working_hour > ?", "09:00:00").Count(&dashboardAttendance.UneligibleWorkingHour)
	result = ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Not("notes = ?", "").Count(&dashboardAttendance.Penugasan)
	return dashboardAttendance, result.Error
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
