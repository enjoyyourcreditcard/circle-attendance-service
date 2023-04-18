package mysql

import (
	"circle/domain"
	"context"
	"time"

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
	var attendances 		[]domain.Attendance
	var dashboardAttendance domain.DashboardAttendance

	layout 				:= "02-01-2006 15:04:05"
	parsedStartAt, err 	:= time.Parse(layout, startAt)
	if err != nil { return dashboardAttendance, err }
	
	parsedEndAt, err := time.Parse(layout, endAt)
	if err != nil { return dashboardAttendance, err }

	duration 		:= parsedEndAt.Sub(parsedStartAt)
	period 			:= int(duration.Hours() / 24) + 1
	weekendCount 	:= 0

	for d := parsedStartAt; d.Before(parsedEndAt); d = d.AddDate(0, 0, 1) {
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday { weekendCount++ }
	}

	totalRecords 	:= period - weekendCount
	query 			:= ar.conn.Model(&attendance).Where("start_at BETWEEN ? AND ? AND user_id = ?", startAt, endAt, userId).Not("DAYOFWEEK(STR_TO_DATE(start_at, '%d-%m-%Y')) IN (1,7)").Session(&gorm.Session{})
	result 			:= query.Find(&attendances)
	if result.Error != nil { return dashboardAttendance, result.Error }
	
	result = query.Count(&dashboardAttendance.WorkingDay)
	result = query.Count(&dashboardAttendance.TotalClockin)
	result = query.Not("end_at LIKE ?", "%00:00:00").Count(&dashboardAttendance.TotalClockout)
	result = query.Where("worktype = ?", "wfh").Count(&dashboardAttendance.TotalWfh)
	result = query.Where("worktype = ?", "wfo").Count(&dashboardAttendance.TotalWfo)
	result = query.Where("TIME(STR_TO_DATE(start_at, '%d-%m-%Y %H:%i:%s')) > ?", "09:15:00").Count(&dashboardAttendance.LateIn)
	result = query.Where("TIME(STR_TO_DATE(start_at, '%d-%m-%Y %H:%i:%s')) < ?", "09:15:00").Count(&dashboardAttendance.EarlyIn)
	result = query.Where("TIME(STR_TO_DATE(end_at, '%d-%m-%Y %H:%i:%s')) < ?", "18:00:00").Not("TIME(STR_TO_DATE(end_at, '%d-%m-%Y %H:%i:%s')) = ?", "00:00:00").Count(&dashboardAttendance.EarlyOut)
	result = query.Where("status_start = ?", "inside_area").Count(&dashboardAttendance.InsideArea)
	result = query.Not("status_start != ?", "inside_area").Count(&dashboardAttendance.OutsideArea)
	result = query.Where("status_start = ?", "inside_other_area").Count(&dashboardAttendance.InsideOtherArea)
	result = query.Where("type = ?", "shifting").Count(&dashboardAttendance.Shifting)
	result = query.Not("type = ?", "shifting").Count(&dashboardAttendance.OfficeHour)
	result = query.Not("working_hour = ?", "").Where("working_hour < ?", "09:00:00").Count(&dashboardAttendance.UneligibleWorkingHour)
	result = query.Not("notes = ?", "").Count(&dashboardAttendance.Penugasan)
	result = query.Where("absen_status = ?", "sick").Count(&dashboardAttendance.Sick)
	result = query.Where("absen_status = ?", "izin").Count(&dashboardAttendance.Izin)

	numRecords 			:= len(attendances)
	remainingRecords 	:= totalRecords - numRecords

	dashboardAttendance.NonWorkingDay = int64(remainingRecords)
	dashboardAttendance.Alpa 		  = int64(remainingRecords)

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