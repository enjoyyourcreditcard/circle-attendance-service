package helper

import (
	"circle/domain"
	"strconv"
)

func ConvertUserArrToUserID(userArr []domain.UserAPIResponse) []int {
	var userIds []int
	for _, user := range userArr {
		userIds = append(userIds, user.ID)
	}
	return userIds
}

func CombineUserAndAttendance(users []domain.UserAPIResponse, attendances []domain.Attendance) []domain.AttendanceExcel {
	attendanceExcel := make([]domain.AttendanceExcel, 0)
	for _, attendance := range attendances {
		var user domain.UserAPIResponse
		user = searchUser(users, attendance.UserId)
		item := domain.AttendanceExcel{
			Nik:      user.Nik,
			Name:     user.Name,
			Email:    user.Email,
			Date:     attendance.StartAt,
			ClockIn:  attendance.StartAt,
			Timezone: attendance.Timezone,
			Worktype: attendance.Worktype,
		}
		attendanceExcel = append(attendanceExcel, item)
	}
	return attendanceExcel
}

func searchUser(users []domain.UserAPIResponse, userId string) domain.UserAPIResponse {
	for _, u := range users {
		if strconv.Itoa(u.ID) == userId {
			return u
			break
		}
	}
	return domain.UserAPIResponse{}
}
