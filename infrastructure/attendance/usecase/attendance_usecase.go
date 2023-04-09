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

func (au attendanceUsecase) GetUserAttendanceMonthly(ctx context.Context, formatedCurrentMY string, userId string) (domain.AttendanceMonthly, error) {
	res, err := au.attendanceRepo.GetUserAttendanceMonthly(ctx, formatedCurrentMY, userId)
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
	latestAttendance, err := au.attendanceRepo.GetUserLastAttendance(ctx, userId)
	if err != nil { return "", err }
	
	attendanceId 	:= latestAttendance.ID
	loc, err 		:= time.LoadLocation("Local")
	if err != nil { return "", err }

	now 				:= time.Now()
	parsedStartAt, err 	:= time.ParseInLocation("02-01-2006 15:04:05", latestAttendance.StartAt, loc)
	if err != nil { return "", err }
	
	totalDuration 				:= int(now.Sub(parsedStartAt).Seconds())
	workingHour 				:= fmt.Sprintf("%02d:%02d:%02d", totalDuration/3600, (totalDuration%3600)/60, totalDuration%60)
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

func (au attendanceUsecase) Hello() string {
	response := au.attendanceRepo.Hello()
	return response
}