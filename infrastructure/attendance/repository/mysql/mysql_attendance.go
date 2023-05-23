package mysql

import (
	"circle/domain"
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type mysqlmysqlAttendanceRepository struct {
	conn *gorm.DB
}

func NewMysqlmysqlAttendanceRepository(conn *gorm.DB) domain.AttendanceMysqlRepository {
	return &mysqlmysqlAttendanceRepository{conn}
}

func (ar mysqlmysqlAttendanceRepository) GetUserLastAttendance(ctx context.Context, userId string) (domain.Attendance, error) {
	var attendance domain.Attendance

	result := ar.conn.Where("user_id = ?", userId).Last(&attendance)
	return attendance, result.Error
}

func (ar mysqlmysqlAttendanceRepository) GetUserDashboardAttendance(ctx context.Context, userId int, startAt string, endAt string) (domain.DashboardAttendance, error) {
	var attendance domain.Attendance
	var attendances []domain.Attendance
	var dashboardAttendance domain.DashboardAttendance

	layout := "02-01-2006"
	parsedStartAt, err := time.Parse(layout, startAt)
	if err != nil {
		return dashboardAttendance, err
	}

	parsedEndAt, err := time.Parse(layout, endAt)
	if err != nil {
		return dashboardAttendance, err
	}

	duration := parsedEndAt.Sub(parsedStartAt)
	period := int(duration.Hours() / 24)

	query := ar.conn.Model(&attendance).Where("STR_TO_DATE(start_at, '%d-%m-%Y') BETWEEN ? AND ? AND user_id = ?", parsedStartAt.Format("2006-01-02"), parsedEndAt.Format("2006-01-02"), userId).Session(&gorm.Session{})
	result := query.Find(&attendances)
	if result.Error != nil {
		return dashboardAttendance, result.Error
	}

	result = query.Count(&dashboardAttendance.WorkingDay)
	result = query.Count(&dashboardAttendance.TotalClockin)
	result = query.Not("end_at LIKE ?", "%00:00:00").Count(&dashboardAttendance.TotalClockout)
	result = query.Where("worktype = ?", "wfh").Count(&dashboardAttendance.TotalWfh)
	result = query.Where("worktype = ?", "wfo").Count(&dashboardAttendance.TotalWfo)
	result = query.Where("reason LIKE ?", "%late_in%").Count(&dashboardAttendance.LateIn)
	result = query.Where("reason LIKE ?", "%early_in%").Count(&dashboardAttendance.EarlyIn)
	result = query.Where("reason LIKE ?", "%early_out%").Count(&dashboardAttendance.EarlyOut)
	result = query.Where("status_start = ?", "inside_area").Count(&dashboardAttendance.InsideArea)
	result = query.Not("status_start != ?", "inside_area").Count(&dashboardAttendance.OutsideArea)
	result = query.Where("status_start = ?", "inside_other_area").Count(&dashboardAttendance.InsideOtherArea)
	result = query.Where("type = ?", "shifting").Count(&dashboardAttendance.Shifting)
	result = query.Not("type = ?", "shifting").Count(&dashboardAttendance.OfficeHour)
	result = query.Not("working_hour = ?", "").Where("working_hour < ?", "09:00:00").Count(&dashboardAttendance.UneligibleWorkingHour)
	result = query.Not("notes = ?", "").Count(&dashboardAttendance.Penugasan)
	result = query.Where("absen_status = ?", "sick").Count(&dashboardAttendance.Sick)
	result = query.Where("absen_status = ?", "izin").Count(&dashboardAttendance.Izin)

	numRecords := len(attendances)
	remainingRecords := period - numRecords

	dashboardAttendance.NonWorkingDay = int64(remainingRecords)
	dashboardAttendance.Alpa = int64(remainingRecords)

	return dashboardAttendance, result.Error
}

func (ar mysqlmysqlAttendanceRepository) GetUserAttendanceData(ctx context.Context, userId string, startAt string, endAt string) ([]domain.Attendance, error) {

	var attendances []domain.Attendance
	query := ar.conn.Model(&attendances)
	if len(userId) != 0 {
		query.Where("user_id = ?", userId)
	}
	//start :=
	start, err := time.Parse("02-01-2006", strings.Split(startAt, " ")[0])
	if err != nil {
		fmt.Println("Error:", err)
	}

	end, err := time.Parse("02-01-2006", strings.Split(endAt, " ")[0])
	if err != nil {
		fmt.Println("Error:", err)
	}
	result := query.Where("STR_TO_DATE(start_at, '%d-%m-%Y') BETWEEN ? AND ?", start.Format("2006-01-02"), end.Format("2006-01-02")).Find(&attendances)
	return attendances, result.Error

}

func (ar mysqlmysqlAttendanceRepository) CheckAbsen(ctx context.Context, userId string, formatedCurrentDate string) (int, error) {
	var counted int64
	var Attendance domain.Attendance

	err := ar.conn.Where("user_id = ?", userId).Where("start_at LIKE ?", formatedCurrentDate+"%").Find(&Attendance).Count(&counted)
	result := int(counted)
	return result, err.Error
}

func (ar mysqlmysqlAttendanceRepository) CreateAbsen(ctx context.Context, attendance *domain.Attendance) (*domain.Attendance, error) {
	result := ar.conn.Create(attendance)
	return attendance, result.Error
}

func (ar mysqlmysqlAttendanceRepository) UpdateAbsen(ctx context.Context, endAttendance *domain.EndAttendance, attendanceId int) error {
	var attendance domain.Attendance

	result := ar.conn.Model(attendance).Where("id = ?", attendanceId).Updates(*endAttendance)
	return result.Error
}

func (ar mysqlmysqlAttendanceRepository) PostAttendanceNotes(ctx context.Context, userId string, notes string) error {
	var attendance domain.Attendance

	result := ar.conn.Where("user_id = ?", userId).Last(&attendance).Update("notes", notes)
	return result.Error
}

func (ar mysqlmysqlAttendanceRepository) GetAttendanceByUserIDAndDateRange(ctx context.Context, userId []int, startAt time.Time, endAt time.Time) ([]domain.Attendance, error) {
	var attendances []domain.Attendance
	query := ar.conn.Model(&attendances)

	if len(userId) != 0 {
		query.Where("user_id IN (?)", userId)
	}
	result := query.Where("STR_TO_DATE(start_at, '%d-%m-%Y') BETWEEN ? AND ?", startAt.Format("2006-01-02"), endAt.Format("2006-01-02")).
		Find(&attendances)
	return attendances, result.Error
}
