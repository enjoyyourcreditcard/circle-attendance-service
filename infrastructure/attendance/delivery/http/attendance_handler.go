package http

import (
	"circle/domain"
	"circle/helper"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
	"time"
)

type AttendanceHandler struct {
	AttendanceUsecase domain.AttendanceUsecase
}

func NewAttendanceHandler(app *fiber.App, atu domain.AttendanceUsecase) {
	handler := &AttendanceHandler{AttendanceUsecase: atu}

	app.Post("start/absen", handler.PostStartAbsen)
	app.Post("end/absen", handler.PostStopAbsen)
}

func (ath *AttendanceHandler) PostStartAbsen(c *fiber.Ctx) error {
	var attendance domain.Attendance
	err := c.BodyParser(&attendance)
	if err != nil { return helper.ResponseIfError(err, c) }

	currentTime 				:= time.Now()
	formatedCurrentDate 		:= currentTime.Format("02-01-2006")
	formatedCurrentTimeDate 	:= currentTime.Format("02-01-2006 15:04:05")
	userId 						:= c.FormValue("user_id")
	attendanceStartAt 			:= formatedCurrentTimeDate
	attendanceEndAt 			:= formatedCurrentTimeDate
	attendanceAbsenStatus 		:= c.FormValue("absen_status")
	attendanceLocationStartId 	:= c.FormValue("location_id")
	attendanceStatusStart 		:= c.FormValue("status_start")
	convertedAttendanceLocationStartId, err := strconv.Atoi(attendanceLocationStartId)
	if err != nil { return helper.ResponseIfError(err, c) }

	attendance.UserId 			= userId
	attendance.StartAt 			= attendanceStartAt
	attendance.EndAt 			= attendanceEndAt
	attendance.AbsenStatus 		= attendanceAbsenStatus
	attendance.LocationStartID 	= convertedAttendanceLocationStartId
	attendance.StatusStart 		= attendanceStatusStart
	data, err := ath.AttendanceUsecase.PostStartAbsen(c.Context(), &attendance, formatedCurrentDate)
	if err != nil { return helper.ResponseIfError(err, c) }

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}

func (ath *AttendanceHandler) PostStopAbsen(c *fiber.Ctx) error {
	var endAttendance domain.EndAttendance
	currentTime 							:= time.Now()
	userId 									:= c.FormValue("user_id")
	attendanceEndAt 						:= currentTime.Format("02-01-2006 15:04:05")
	attendanceReason 						:= c.FormValue("reason")
	attendanceLocationEndId 				:= c.FormValue("location_id")
	attendanceLatLongPositionEnd 			:= c.FormValue("lat_long_position_end")
	attendanceStatusEnd 					:= c.FormValue("status_end")
	convertedAttendanceLocationEndId, err 	:= strconv.Atoi(attendanceLocationEndId)
	if err != nil { return helper.ResponseIfError(err, c) }

	endAttendance.EndAt 				= attendanceEndAt
	endAttendance.Reason 				= attendanceReason
	endAttendance.LocationEndID 		= convertedAttendanceLocationEndId
	endAttendance.LatLongPositionEnd 	= attendanceLatLongPositionEnd
	endAttendance.StatusEnd 			= attendanceStatusEnd

	data, err := ath.AttendanceUsecase.PostStopAbsen(c.Context(), &endAttendance, userId)
	if err != nil { return helper.ResponseIfError(err, c) }

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}
