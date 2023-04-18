package helper

import (
	"circle/domain"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)


func GetChildren(c *fiber.Ctx, parentId string) ([]domain.User, error) {
	var apiResponse 	domain.ApiResponse
	var children 		[]domain.User

	client := new(http.Client)
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://hris.mncplay.id/user/api/heirarky/%s", parentId), nil)
	if err != nil { return children, err }

	res, err := client.Do(req)
	if err != nil { return children, err }
	
	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil { return children, err }

	err = json.Unmarshal(responseBody, &apiResponse)
	if err != nil { return children, err }

	for _, s := range apiResponse.Data.Child { children = append(children, s) }

	return children, err
}