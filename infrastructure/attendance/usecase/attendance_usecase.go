package usecase

import (
	"circle/domain"
	"context"
	"fmt"
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

func (au attendanceUsecase) CheckAbsen(ctx context.Context, userId string, formatedCurrentDate string) (int, error) {
	res, err := au.attendanceRepo.CheckAbsen(ctx, userId, formatedCurrentDate)
	return res, err
}

func (au attendanceUsecase) PostStartAbsen(ctx context.Context, attendance *domain.Attendance, formatedCurrentDate string) (string, error) {
	var res string
	var err error
	var newAttendance *domain.Attendance

	checkAbsen, err := au.attendanceRepo.CheckAbsen(ctx, attendance.UserId, formatedCurrentDate)
	if err != nil { return "", err }

	if checkAbsen > 0 {
		res 				= "Anda sudah absen"
	} else {
		newAttendance, err 	= au.attendanceRepo.CreateAbsen(ctx, attendance)
		res 				= newAttendance.StartAt
	}
	return res, err
}

func (au attendanceUsecase) PostStopAbsen(ctx context.Context, endAttendance *domain.EndAttendance, userId string) (string, error) {
	var attendance domain.Attendance
	
	latestAttendance, err := au.attendanceRepo.GetLatestUserAbsen(ctx, attendance, userId)
	if err != nil { return "Attendance not found", err }
	
	attendanceId 	:= latestAttendance.ID
	loc, err 		:= time.LoadLocation("Local")
	if err != nil { return "", err }

	now 				:= time.Now()
	parsedStartAt, err 	:= time.ParseInLocation("02-01-2006 15:04:05", latestAttendance.StartAt, loc)
	if err != nil { return fmt.Sprintf("There's a problem while parsing %s", latestAttendance.StartAt), err }
	
	totalDuration 				:= int(now.Sub(parsedStartAt).Seconds())
	workingHour 				:= fmt.Sprintf("%02d:%02d:%02d", totalDuration/3600, (totalDuration%3600)/60, totalDuration%60)
	endAttendance.WorkingHour 	= workingHour
	err 						= au.attendanceRepo.UpdateAbsen(ctx, endAttendance, attendanceId)
	if err != nil { return "", err }

	res := endAttendance.EndAt

	return res, err
}

func (au attendanceUsecase) Hello() string {
	response := au.attendanceRepo.Hello()
	return response
}
