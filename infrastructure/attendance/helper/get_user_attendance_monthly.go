package helper

import (
	"circle/domain"
)

func GetUserAttendanceMonthly(attendanceMonthly domain.AttendanceMonthly, userId string, totalAbsen int64, totalWfo int64, totalWfh int64, totalOff int64) domain.AttendanceMonthly {
	attendanceMonthly.UserId 		= userId
	attendanceMonthly.TotalAbsen 	= totalAbsen
	attendanceMonthly.TotalWfo 		= totalWfo
	attendanceMonthly.TotalWfh 		= totalWfh
	attendanceMonthly.TotalOff 		= totalOff

	return attendanceMonthly
}