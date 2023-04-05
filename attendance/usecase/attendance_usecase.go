package usecase

import (
	"circle/domain"
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

func PostStartAbsen()  {
	
}