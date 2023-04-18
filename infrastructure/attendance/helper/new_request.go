package helper

import (
	"circle/domain"
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

func GetChildren(c *fiber.Ctx, parentId string) ([]domain.User, error) {
	var apiResponse 	domain.ApiResponse
	var children 		[]domain.User

	client := resty.New()
    resp, err := client.R().
        Get("https://hris.mncplay.id/user/api/heirarky/" + parentId)
	if err != nil { return children, err }
	
    err = json.Unmarshal(resp.Body(), &apiResponse)
	if err != nil { return children, err }

	data := apiResponse.Data.Child
	for _, s := range data { children = append(children, s) }
	
	return children, err
}