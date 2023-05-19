package usecase

import (
	"circle/domain"
	"circle/infrastructure/attendance/usecase/helper"
	"context"
	"fmt"
	"time"
)

type attendanceUsecase struct {
	mysqlAttendanceRepo domain.AttendanceMysqlRepository
	apiAttendanceRepo   domain.AttendanceAPIRepository

	contextTimeout time.Duration
}

func NewAttendanceUsecase(m domain.AttendanceMysqlRepository, a domain.AttendanceAPIRepository, timeout time.Duration) domain.AttendanceUsecase {
	return &attendanceUsecase{
		mysqlAttendanceRepo: m,
		apiAttendanceRepo:   a,
		contextTimeout:      timeout,
	}
}

func (au attendanceUsecase) GetUserLastAttendance(ctx context.Context, userId string) (domain.Attendance, error) {
	res, err := au.mysqlAttendanceRepo.GetUserLastAttendance(ctx, userId)
	if err != nil {
		res = domain.Attendance{}
		res.UserId = userId
		res.Worktype = "Belum absen"
		return res, nil
	}

	return res, err
}

func (au attendanceUsecase) GetUserDashboardAttendance(ctx context.Context, userId int, startAt string, endAt string) (domain.DashboardAttendance, error) {
	res, err := au.mysqlAttendanceRepo.GetUserDashboardAttendance(ctx, userId, startAt, endAt)
	return res, err
}

func (au attendanceUsecase) GetUserAttendanceData(ctx context.Context, userId string, startAt string, endAt string) ([]domain.Attendance, error) {
	res, err := au.mysqlAttendanceRepo.GetUserAttendanceData(ctx, userId, startAt, endAt)
	return res, err
}

func (au attendanceUsecase) GetChildDashboardAttendance(ctx context.Context, startAt string, endAt string, parentID int) (domain.DashboardAttendanceChildren, error) {
	var dashboardChildren []domain.DashboardAttendance
	var dashboard domain.DashboardAttendanceChildren
	children, err := helper.GetChildren(parentID)

	if err != nil {
		return domain.DashboardAttendanceChildren{}, err
	}

	for _, child := range children {
		userId := child.ID
		res, err := au.mysqlAttendanceRepo.GetUserDashboardAttendance(ctx, userId, startAt, endAt)
		if err != nil {
			fmt.Println(err)
			//return domain.DashboardAttendance{}, err
		}

		dashboardChildren = append(dashboardChildren, res)
	}

	for _, data := range dashboardChildren {
		dashboard.WorkingDay += data.WorkingDay
		dashboard.NonWorkingDay += data.NonWorkingDay
		dashboard.Holiday += data.Holiday
		dashboard.TotalClockin += data.TotalClockin
		dashboard.TotalClockout += data.TotalClockout
		dashboard.TotalWfh += data.TotalWfh
		dashboard.TotalWfo += data.TotalWfo
		dashboard.LateIn += data.LateIn
		dashboard.EarlyIn += data.EarlyIn
		dashboard.EarlyOut += data.EarlyOut
		dashboard.InsideArea += data.InsideArea
		dashboard.OutsideArea += data.OutsideArea
		dashboard.InsideOtherArea += data.InsideOtherArea
		dashboard.Shifting += data.Shifting
		dashboard.OfficeHour += data.OfficeHour
		dashboard.Alpa += data.Alpa
		dashboard.Sick += data.Sick
		dashboard.Izin += data.Izin
		dashboard.Leave += data.Leave
		dashboard.UneligibleWorkingHour += data.UneligibleWorkingHour
		dashboard.Penugasan += data.Penugasan
	}
	dashboard.UserData = children

	return dashboard, err
}

func (au attendanceUsecase) GetChildDashboardAttendanceDetail(ctx context.Context, startAt string, endAt string, parentID int) ([]domain.DashboardAttendanceChildren, error) {

	children, err := helper.GetChildren(parentID)
	if err != nil {
		return nil, err
	}

	//var res domain.DashboardAttendance
	var dashboards []domain.DashboardAttendanceChildren
	//var err error

	for _, child := range children {
		userId := child.ID
		res, _ := au.mysqlAttendanceRepo.GetUserDashboardAttendance(ctx, userId, startAt, endAt)

		dashboard := domain.DashboardAttendanceChildren{
			WorkingDay:            res.WorkingDay,
			NonWorkingDay:         res.NonWorkingDay,
			Holiday:               res.Holiday,
			TotalClockin:          res.TotalClockin,
			TotalClockout:         res.TotalClockout,
			TotalWfh:              res.TotalWfh,
			TotalWfo:              res.TotalWfo,
			LateIn:                res.LateIn,
			EarlyIn:               res.EarlyIn,
			EarlyOut:              res.EarlyOut,
			InsideArea:            res.InsideArea,
			OutsideArea:           res.OutsideArea,
			InsideOtherArea:       res.InsideOtherArea,
			Shifting:              res.Shifting,
			OfficeHour:            res.OfficeHour,
			Alpa:                  res.Alpa,
			Sick:                  res.Sick,
			Izin:                  res.Izin,
			Leave:                 res.Leave,
			UneligibleWorkingHour: res.UneligibleWorkingHour,
			Penugasan:             res.Penugasan,
			UserData:              []domain.User{child},
		}
		dashboards = append(dashboards, dashboard)

	}

	return dashboards, nil
}

func (au attendanceUsecase) PostClockIn(ctx context.Context, attendance *domain.Attendance, formatedCurrentDate string) (string, error) {
	var res string
	var err error
	var newAttendance *domain.Attendance

	checkAbsen, err := au.mysqlAttendanceRepo.CheckAbsen(ctx, attendance.UserId, formatedCurrentDate)
	if err != nil {
		return "", err
	}

	if checkAbsen > 0 {
		res = "Anda sudah clock in"
	} else {
		newAttendance, err = au.mysqlAttendanceRepo.CreateAbsen(ctx, attendance)
		res = newAttendance.StartAt
	}
	return res, err
}

func (au attendanceUsecase) PostClockOut(ctx context.Context, endAttendance *domain.EndAttendance, userId string) (string, error) {
	latestAttendance, err := au.mysqlAttendanceRepo.GetUserLastAttendance(ctx, userId)
	if err != nil {
		return "", err
	}

	if latestAttendance.WorkingHour != "" {
		return "Anda sudah clock out", nil
	}

	attendanceId := latestAttendance.ID
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return "", err
	}

	parsedStartAt, err := time.ParseInLocation("02-01-2006 15:04:05", latestAttendance.StartAt, loc)
	if err != nil {
		return "", err
	}

	parsedEndAt, err := time.ParseInLocation("02-01-2006 15:04:05", endAttendance.EndAt, loc)
	if err != nil {
		return "", err
	}

	duration := parsedEndAt.Sub(parsedStartAt)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	workingHour := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	endAttendance.WorkingHour = workingHour
	err = au.mysqlAttendanceRepo.UpdateAbsen(ctx, endAttendance, attendanceId)
	if err != nil {
		return "", err
	}

	res := endAttendance.EndAt

	return res, err
}

func (au attendanceUsecase) PostAttendanceNotes(ctx context.Context, userId string, notes string) (string, error) {
	err := au.mysqlAttendanceRepo.PostAttendanceNotes(ctx, userId, notes)
	if err != nil {
		return "", err
	}

	res := "Notes berhasil ditambahkan"
	return res, err
}

func (au attendanceUsecase) HandleAttendanceExcel(ctx context.Context, name string, nik string, unit_bisnis string, status_karyawan string, regional string, startAt time.Time, endAt time.Time) ([]domain.AttendanceExcel, error) {

	user, err := au.apiAttendanceRepo.Find(ctx, name, nik, unit_bisnis, status_karyawan, regional)
	userId := helper.ConvertUserArrToUserID(user)
	attendances, err := au.mysqlAttendanceRepo.GetAttendanceByUserIDAndDateRange(ctx, userId, startAt, endAt)
	attendanceExcel := helper.CombineUserAndAttendance(user, attendances)

	return attendanceExcel, err
}

func (au attendanceUsecase) GetAttendanceByUserID(ctx context.Context, startAt time.Time, endAt time.Time, userId []int) ([]domain.Attendance, error) {
	attendances, err := au.mysqlAttendanceRepo.GetAttendanceByUserIDAndDateRange(ctx, userId, startAt, endAt)
	return attendances, err
}
