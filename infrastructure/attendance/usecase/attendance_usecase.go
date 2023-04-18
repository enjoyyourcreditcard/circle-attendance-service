package usecase

import (
	"circle/domain"
	"context"
	"fmt"
	"strconv"
	"time"
)

type attendanceUsecase struct {
	attendanceRepo domain.AttendanceRepository
	contextTimeout time.Duration
}

func NewAttendanceUsecase(a domain.AttendanceRepository, timeout time.Duration) domain.AttendanceUsecase {
	return &attendanceUsecase{
		attendanceRepo: a,
		contextTimeout: timeout,
	}
}

func (au attendanceUsecase) GetUserLastAttendance(ctx context.Context, userId string) (domain.Attendance, error)  {
	res, err := au.attendanceRepo.GetUserLastAttendance(ctx, userId)
	if err != nil {
		res 			= domain.Attendance{}
		res.UserId 		= userId
		res.Worktype	= "Belum absen"
		return res, nil
	}

	return res, err
}

func (au attendanceUsecase) GetUserDashboardAttendance(ctx context.Context, userId string, startAt string, endAt string) (domain.DashboardAttendance, error) {
	res, err 		:= au.attendanceRepo.GetUserDashboardAttendance(ctx, userId, startAt, endAt)
	return res, err
}

func (au attendanceUsecase) GetUserAttendanceData(ctx context.Context, userId string, startAt string, endAt string) ([]domain.Attendance, error) {
	res, err := au.attendanceRepo.GetUserAttendanceData(ctx, userId, startAt, endAt)
	return res, err
}

func (au attendanceUsecase) GetChildDashboardAttendance(ctx context.Context, startAt string, endAt string, children []domain.User) (domain.DashboardAttendance, error) {
	var res 		domain.DashboardAttendance
	var dashboard 	domain.DashboardAttendance
	var dashboards 	[]domain.DashboardAttendance
	var err 		error
	
	for _, child := range children {
		userId 			:= strconv.Itoa(child.ID)
		res, err 		= au.attendanceRepo.GetUserDashboardAttendance(ctx, userId, startAt, endAt)
		if err != nil { return domain.DashboardAttendance{}, err }

		dashboards 		= append(dashboards, res)
	}

	for _, data := range dashboards {
        dashboard.WorkingDay 			+= data.WorkingDay
        dashboard.NonWorkingDay 		+= data.NonWorkingDay
        dashboard.Holiday 				+= data.Holiday
        dashboard.TotalClockin 			+= data.TotalClockin
        dashboard.TotalClockout 		+= data.TotalClockout
        dashboard.TotalWfh 				+= data.TotalWfh
        dashboard.TotalWfo 				+= data.TotalWfo
        dashboard.LateIn 				+= data.LateIn
        dashboard.EarlyIn 				+= data.EarlyIn
        dashboard.EarlyOut 				+= data.EarlyOut
        dashboard.InsideArea 			+= data.InsideArea
        dashboard.OutsideArea 			+= data.OutsideArea
        dashboard.InsideOtherArea 		+= data.InsideOtherArea
        dashboard.Shifting 				+= data.Shifting
        dashboard.OfficeHour 			+= data.OfficeHour
        dashboard.Alpa 					+= data.Alpa
        dashboard.Sick 					+= data.Sick
        dashboard.Izin 					+= data.Izin
        dashboard.Leave 				+= data.Leave
        dashboard.UneligibleWorkingHour += data.UneligibleWorkingHour
        dashboard.Penugasan 			+= data.Penugasan
    }

	return dashboard, err
}

func (au attendanceUsecase) GetChildDashboardAttendanceDetail(ctx context.Context, startAt string, endAt string, children []domain.User) ([]domain.DashboardAttendance, error) {
	var res 		domain.DashboardAttendance
	var dashboards 	[]domain.DashboardAttendance
	var err 		error
	
	for _, child := range children {
		userId 			:= strconv.Itoa(child.ID)
		res, err 		= au.attendanceRepo.GetUserDashboardAttendance(ctx, userId, startAt, endAt)
		res.UserData 	= child
		dashboards		= append(dashboards, res)
	}
	
	return dashboards, err
}

func (au attendanceUsecase) PostClockIn(ctx context.Context, attendance *domain.Attendance, formatedCurrentDate string) (string, error) {
	var res string
	var err error
	var newAttendance *domain.Attendance

	checkAbsen, err := au.attendanceRepo.CheckAbsen(ctx, attendance.UserId, formatedCurrentDate)
	if err != nil { return "", err }

	if checkAbsen > 0 {
		res 				= "Anda sudah clock in"
	} else {
		newAttendance, err 	= au.attendanceRepo.CreateAbsen(ctx, attendance)
		res 				= newAttendance.StartAt
	}
	return res, err
}

func (au attendanceUsecase) PostClockOut(ctx context.Context, endAttendance *domain.EndAttendance, userId string) (string, error) {
	latestAttendance, err := au.attendanceRepo.GetUserLastAttendance(ctx, userId)
	if err != nil { return "", err }

	if latestAttendance.WorkingHour != "" { return "Anda sudah clock out", nil }
	
	attendanceId 	:= latestAttendance.ID
	loc, err 		:= time.LoadLocation("Local")
	if err != nil { return "", err }

	parsedStartAt, err 	:= time.ParseInLocation("02-01-2006 15:04:05", latestAttendance.StartAt, loc)
	if err != nil { return "", err }

	parsedEndAt, err 	:= time.ParseInLocation("02-01-2006 15:04:05", endAttendance.EndAt, loc)
	if err != nil { return "", err }
	
	duration 	:= parsedEndAt.Sub(parsedStartAt)
    hours 		:= int(duration.Hours())
    minutes 	:= int(duration.Minutes()) % 60
    seconds 	:= int(duration.Seconds()) % 60
    workingHour := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	endAttendance.WorkingHour 	= workingHour
	err 						= au.attendanceRepo.UpdateAbsen(ctx, endAttendance, attendanceId)
	if err != nil { return "", err }

	res := endAttendance.EndAt

	return res, err
}

func (au attendanceUsecase) PostAttendanceNotes(ctx context.Context, userId string, notes string) (string, error) {
	err := au.attendanceRepo.PostAttendanceNotes(ctx, userId, notes)
	if err != nil { return "", err }

	res := "Notes berhasil ditambahkan"
	return res, err
}