package helper

import (
	"circle/domain"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func ResponseIfError(err error, c *fiber.Ctx) error {
	if err != nil {
		return c.JSON(domain.WebResponse{Status: http.StatusInternalServerError, Data: nil, Message: err.Error()})
	} else {
		return nil
	}
}