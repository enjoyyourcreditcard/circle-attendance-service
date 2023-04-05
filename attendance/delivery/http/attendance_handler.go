package http

import (
	"circle/domain"
	"github.com/gofiber/fiber/v2"
)

type AttendanceHandler struct {
	AttendanceUsecase domain.AttendanceUsecase
}

func NewAttendanceHandler(app *fiber.App, us domain.AttendanceUsecase) {
	handler := &AttendanceHandler{ AttendanceUsecase: us }

	app.Post("DEV/mobile/api/start/absen", handler.PostStartAbsen)
}

func (u *AttendanceHandler) PostStartAbsen(c *fiber.Ctx) error {

	return nil
}
