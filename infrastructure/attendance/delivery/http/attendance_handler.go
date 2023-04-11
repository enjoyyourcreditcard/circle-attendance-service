package http

import (
	"circle/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AttendanceHandler struct {
	AttendanceUsecase domain.AttendanceUsecase
}

func NewAttendanceHandler(app *fiber.App, atu domain.AttendanceUsecase) {
	handler := &AttendanceHandler{AttendanceUsecase: atu}

	app.Get("absen/user/:user_id", handler.GetUserLastAttendance)
	app.Get("absen/user/bulanan/:user_id", handler.GetUserAttendanceMonthly)
	
	app.Post("start/absen", handler.PostClockIn)
	app.Post("end/absen", handler.PostClockOut)
	app.Post("attendance/notes", handler.PostAttendanceNotes)
}

func (ath *AttendanceHandler) GetUserLastAttendance(c *fiber.Ctx) error  {
	userId 		:= c.Params("user_id")
	data, err	:= ath.AttendanceUsecase.GetUserLastAttendance(c.Context(), userId)
	if err != nil { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil}) }
	
	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}

func (ath *AttendanceHandler) GetUserAttendanceMonthly(c *fiber.Ctx) error {
	userId 				:= c.Params("user_id")
	currentTime			:= time.Now()
	formatedCurrentMY 	:= currentTime.Format("01-2006")
	data, err 			:= ath.AttendanceUsecase.GetUserAttendanceMonthly(c.Context(), formatedCurrentMY, userId)
	if err != nil { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil}) }

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}

func (ath *AttendanceHandler) PostClockIn(c *fiber.Ctx) error {
	var attendance 				domain.Attendance
	var currentTime 			time.Time
	var formatedCurrentDate 	string

	err := c.BodyParser(&attendance)
	if err != nil { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil}) }

	timeBody := c.FormValue("time")
	if timeBody != "" {
		t, err := time.Parse("02-01-2006 15:04:05", timeBody)
		if err != nil { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil}) }

		formatedCurrentDate = t.Format("02-01-2006")
	} else {
		currentTime 			= time.Now()
		formatedCurrentDate 	= currentTime.Format("02-01-2006")
		timeBody 				= currentTime.Format("02-01-2006 15:04:05")
	}

	userId 						:= c.FormValue("user_id")
	attendanceStartAt 			:= timeBody
	attendanceEndAt 			:= timeBody
	attendanceAbsenStatus 		:= c.FormValue("absen_status")
	attendanceLocationStartId 	:= c.FormValue("location_id")
	attendanceStatusStart 		:= c.FormValue("status_start")
	latLongPositionStart		:= c.FormValue("lat_long_position_start")
	convertedAttendanceLocationStartId, err := strconv.Atoi(attendanceLocationStartId)
	if err != nil { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil}) }

	attendance.UserId 				= userId
	attendance.StartAt 				= attendanceStartAt
	attendance.EndAt 				= attendanceEndAt
	attendance.AbsenStatus 			= attendanceAbsenStatus
	attendance.LocationStartID 		= convertedAttendanceLocationStartId
	attendance.StatusStart 			= attendanceStatusStart
	attendance.LatLongPositionStart	= latLongPositionStart
	data, err := ath.AttendanceUsecase.PostClockIn(c.Context(), &attendance, formatedCurrentDate)
	if err != nil { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil}) }
	if data == "Anda sudah clock in" { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: data, Data: nil}) }

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}

func (ath *AttendanceHandler) PostClockOut(c *fiber.Ctx) error {
	var endAttendance 			domain.EndAttendance
	var currentTime				time.Time

	timeBody := c.FormValue("time")
	if timeBody == "" {
		currentTime	= time.Now()
		timeBody 	= currentTime.Format("02-01-2006 15:04:05")
	}
	
	userId 									:= c.FormValue("user_id")
	attendanceReason 						:= c.FormValue("reason")
	attendanceLocationEndId 				:= c.FormValue("location_id")
	attendanceLatLongPositionEnd 			:= c.FormValue("lat_long_position_end")
	attendanceStatusEnd 					:= c.FormValue("status_end")
	convertedAttendanceLocationEndId, err 	:= strconv.Atoi(attendanceLocationEndId)
	if err != nil { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil}) }

	endAttendance.EndAt 				= timeBody
	endAttendance.Reason 				= attendanceReason
	endAttendance.LocationEndID 		= convertedAttendanceLocationEndId
	endAttendance.LatLongPositionEnd 	= attendanceLatLongPositionEnd
	endAttendance.StatusEnd 			= attendanceStatusEnd

	data, err := ath.AttendanceUsecase.PostClockOut(c.Context(), &endAttendance, userId)
	if err != nil { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil}) }
	
	if data == "Anda sudah clock out" { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: data, Data: nil}) }

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}

func (ath *AttendanceHandler) PostAttendanceNotes(c *fiber.Ctx) error {
	userId 		:= c.FormValue("user_id")
	// regional 	:= c.FormValue("regional")
	notes 		:= c.FormValue("notes")
	data, err 	:= ath.AttendanceUsecase.PostAttendanceNotes(c.Context(), userId, notes)
	if err != nil { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil}) }

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}