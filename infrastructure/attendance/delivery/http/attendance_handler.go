package http

import (
	"bytes"
	"circle/domain"
	"fmt"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AttendanceHandler struct {
	AttendanceUsecase domain.AttendanceUsecase
}

func NewAttendanceHandler(app *fiber.App, atu domain.AttendanceUsecase) {
	handler := &AttendanceHandler{AttendanceUsecase: atu}

	app.Get("absen/user/:user_id", handler.GetUserLastAttendance)
	app.Get("absen/user/:user_id/:start_at/:end_at", handler.GetUserAttendanceData)

	app.Get("absen/user/dashboard/:user_id/:start_at/:end_at", handler.GetUserDashboardAttendance)
	app.Get("absen/user/child/dashboard/:parent_id/:start_at/:end_at", handler.GetChildDashboardAttendance)
	app.Get("absen/user/child/array/dashboard/:parent_id/:start_at/:end_at", handler.GetChildDashboardAttendanceDetail)

	app.Post("start/absen", handler.PostClockIn)
	app.Post("end/absen", handler.PostClockOut)
	app.Post("attendance/notes", handler.PostAttendanceNotes)

	app.Get("absen/excel", handler.handleAbsenExcel)
	app.Get("absen/user", handler.GetAbsenUser)

}

func (ath *AttendanceHandler) GetUserDashboardAttendance(c *fiber.Ctx) error {
	userId := c.Params("user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	startAt := c.Params("start_at")
	endAt := c.Params("end_at")
	startAtTime := startAt + " 00:00:01"
	endAtTime := endAt + " 23:59:59"
	data, err := ath.AttendanceUsecase.GetUserDashboardAttendance(c.Context(), userIdInt, startAtTime, endAtTime)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}

func (ath *AttendanceHandler) GetUserLastAttendance(c *fiber.Ctx) error {
	userId := c.Params("user_id")
	data, err := ath.AttendanceUsecase.GetUserLastAttendance(c.Context(), userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}

func (ath *AttendanceHandler) GetUserAttendanceData(c *fiber.Ctx) error {
	userId := c.Params("user_id")
	startAt := c.Params("start_at")
	endAt := c.Params("end_at")
	startAtTime := startAt + " 00:00:01"
	endAtTime := endAt + " 23:59:59"
	data, err := ath.AttendanceUsecase.GetUserAttendanceData(c.Context(), userId, startAtTime, endAtTime)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}

func (ath *AttendanceHandler) GetChildDashboardAttendance(c *fiber.Ctx) error {
	parentIdStr := c.Params("parent_id")
	startAt := c.Params("start_at")
	endAt := c.Params("end_at")

	parentIdInt, err := strconv.Atoi(parentIdStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	//startAtTime := startAt + " 00:00:01"
	//endAtTime := endAt + " 23:59:59"
	//children, err := helper.GetChildren(c, parentId)
	//if err != nil {
	//	return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	//}

	data, err := ath.AttendanceUsecase.GetChildDashboardAttendance(c.Context(), startAt, endAt, parentIdInt)
	//if err != nil {
	//	return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	//}

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}

func (ath *AttendanceHandler) GetChildDashboardAttendanceDetail(c *fiber.Ctx) error {
	parentIdStr := c.Params("parent_id")
	parentIdInt, err := strconv.Atoi(parentIdStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}
	startAt := c.Params("start_at")
	endAt := c.Params("end_at")
	startAtTime := startAt + " 00:00:01"
	endAtTime := endAt + " 23:59:59"

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	data, err := ath.AttendanceUsecase.GetChildDashboardAttendanceDetail(c.Context(), startAtTime, endAtTime, parentIdInt)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}

func (ath *AttendanceHandler) PostClockIn(c *fiber.Ctx) error {
	var attendance domain.Attendance
	var currentTime time.Time
	var formatedCurrentDate string

	err := c.BodyParser(&attendance)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	timeBody := c.FormValue("time")
	if timeBody != "" {
		t, err := time.Parse("02-01-2006 15:04:05", timeBody)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
		}

		formatedCurrentDate = t.Format("02-01-2006")
	} else {
		currentTime = time.Now()
		formatedCurrentDate = currentTime.Format("02-01-2006")
		timeBody = currentTime.Format("02-01-2006 15:04:05")
	}

	timeEndAt := formatedCurrentDate + " 00:00:00"
	userId := c.FormValue("user_id")
	attendanceStartAt := timeBody
	attendanceEndAt := timeEndAt
	attendanceAbsenStatus := c.FormValue("absen_status")
	attendanceLocationStartId := c.FormValue("location_id")
	attendanceStatusStart := c.FormValue("status_start")
	latLongPositionStart := c.FormValue("lat_long_position_start")
	convertedAttendanceLocationStartId, err := strconv.Atoi(attendanceLocationStartId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	attendance.UserId = userId
	attendance.StartAt = attendanceStartAt
	attendance.EndAt = attendanceEndAt
	attendance.AbsenStatus = attendanceAbsenStatus
	attendance.LocationStartID = convertedAttendanceLocationStartId
	attendance.StatusStart = attendanceStatusStart
	attendance.LatLongPositionStart = latLongPositionStart
	data, err := ath.AttendanceUsecase.PostClockIn(c.Context(), &attendance, formatedCurrentDate)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}
	if data == "Anda sudah clock in" {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: data, Data: nil})
	}

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}

func (ath *AttendanceHandler) PostClockOut(c *fiber.Ctx) error {
	var endAttendance domain.EndAttendance
	var currentTime time.Time

	timeBody := c.FormValue("time")
	if timeBody == "" {
		currentTime = time.Now()
		timeBody = currentTime.Format("02-01-2006 15:04:05")
	}

	userId := c.FormValue("user_id")
	attendanceReason := c.FormValue("reason")
	attendanceLocationEndId := c.FormValue("location_id")
	attendanceLatLongPositionEnd := c.FormValue("lat_long_position_end")
	attendanceStatusEnd := c.FormValue("status_end")
	convertedAttendanceLocationEndId, err := strconv.Atoi(attendanceLocationEndId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	endAttendance.EndAt = timeBody
	endAttendance.Reason = attendanceReason
	endAttendance.LocationEndID = convertedAttendanceLocationEndId
	endAttendance.LatLongPositionEnd = attendanceLatLongPositionEnd
	endAttendance.StatusEnd = attendanceStatusEnd

	data, err := ath.AttendanceUsecase.PostClockOut(c.Context(), &endAttendance, userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	if data == "Anda sudah clock out" {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: data, Data: nil})
	}

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}

func (ath *AttendanceHandler) PostAttendanceNotes(c *fiber.Ctx) error {
	userId := c.FormValue("user_id")
	// regional 	:= c.FormValue("regional")
	notes := c.FormValue("notes")
	data, err := ath.AttendanceUsecase.PostAttendanceNotes(c.Context(), userId, notes)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}

func (ath *AttendanceHandler) handleAbsenExcel(c *fiber.Ctx) error {

	startAtStr := c.Query("start_at")
	endAtStr := c.Query("end_at")
	layout := "02-01-2006"

	startAt, err := time.Parse(layout, startAtStr)
	if err != nil {
		fmt.Println(err)
	}

	endAt, err := time.Parse(layout, endAtStr)
	if err != nil {
		fmt.Println(err)
	}

	name := c.Query("name")
	nik := c.Query("nik")
	unitBisnis := c.Query("unit_bisnis")
	statusKaryawan := c.Query("status_karyawan")
	regional := c.Query("regional")

	data, err := ath.AttendanceUsecase.HandleAttendanceExcel(c.Context(), name, nik, unitBisnis, statusKaryawan, regional, startAt, endAt)

	file := excelize.NewFile()
	sheet := "Sheet1"

	file.SetCellValue(sheet, "A1", "nik")
	file.SetCellValue(sheet, "B1", "name")
	file.SetCellValue(sheet, "C1", "email")
	file.SetCellValue(sheet, "D1", "date")
	file.SetCellValue(sheet, "E1", "clock_in")
	file.SetCellValue(sheet, "F1", "worktype")
	file.SetCellValue(sheet, "G1", "timezone")

	for i, each := range data {
		file.SetCellValue(sheet, fmt.Sprintf("A%d", i+2), each.Nik)
		file.SetCellValue(sheet, fmt.Sprintf("B%d", i+2), each.Name)
		file.SetCellValue(sheet, fmt.Sprintf("C%d", i+2), each.Email)
		file.SetCellValue(sheet, fmt.Sprintf("D%d", i+2), each.Date)
		file.SetCellValue(sheet, fmt.Sprintf("E%d", i+2), each.ClockIn)
		file.SetCellValue(sheet, fmt.Sprintf("F%d", i+2), each.Worktype)
		file.SetCellValue(sheet, fmt.Sprintf("G%d", i+2), each.Timezone)
	}

	buffer := bytes.Buffer{}
	if _, err := file.WriteTo(&buffer); err != nil {
		return err
	}

	response := c.Response()
	response.Header.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	response.Header.Set("Content-Disposition", "attachment; filename=people.xlsx")

	if _, err := buffer.WriteTo(response.BodyWriter()); err != nil {
		return err
	}

	return nil
}

func (ath *AttendanceHandler) GetAbsenUser(c *fiber.Ctx) error {
	//validate := validator.New()

	startAtStr := c.Query("start_at")
	if startAtStr == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(domain.WebResponse{Status: 500, Message: "start_at is required", Data: nil})
	}

	endAtStr := c.Query("end_at")
	if endAtStr == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(domain.WebResponse{Status: 500, Message: "end_at is required", Data: nil})
	}

	layout := "02-01-2006"
	userId := c.Query("user_id")
	var arr []int
	if userId != "" {
		strArr := strings.Split(userId, ",")
		arr = make([]int, len(strArr))
		for i, s := range strArr {
			num, err := strconv.Atoi(s)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			arr[i] = num
		}
	}

	startAt, err := time.Parse(layout, startAtStr)
	if err != nil {
		fmt.Println(err)
	}

	endAt, err := time.Parse(layout, endAtStr)
	if err != nil {
		fmt.Println(err)
	}

	attendance, err := ath.AttendanceUsecase.GetAttendanceByUserID(c.Context(), startAt, endAt, arr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}
	return c.Status(200).JSON(domain.WebResponse{Status: 200, Message: "SUCCESS", Data: attendance})
}
