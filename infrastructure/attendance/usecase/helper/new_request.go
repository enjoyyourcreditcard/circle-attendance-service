package helper

import (
	"circle/domain"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"strconv"
)

type HirarkyApiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Child []User `json:"child"`
		//Parent []User `json:"parent"`
	} `json:"data"`
}

type User struct {
	ID             int    `json:"id"`
	ParentID       string `json:"parent_id"`
	Image          string `json:"image"`
	UnitBisnis     string `json:"unit_bisnis"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Roles          string `json:"roles"`
	StatusKaryawan string `json:"status_karyawan"`
	Nik            string `json:"nik"`
	Regional       string `json:"regional"`
	Location       string `json:"location"`
	Status         string `json:"status"`
	Salescode      string `json:"salescode"`
	JoinDate       string `json:"join_date"`
	TglLahir       string `json:"tgl_lahir"`
	Directorate    string `json:"directorate"`
	Department     string `json:"department"`
	JenisKelamin   string `json:"jenis_kelamin"`
}

func GetChildren(parentID int) ([]domain.User, error) {

	var apiResponse HirarkyApiResponse
	var children []domain.User
	var baseUrl string

	if viper.GetString(`server.name`) == "local" {
		baseUrl = "http://localhost"
	} else {
		baseUrl = "https://istudio.mncplay.id"
	}
	parentIDStr := strconv.Itoa(parentID)
	client := resty.New()
	resp, err := client.R().
		Get(baseUrl + "/user/api/heirarky/" + parentIDStr)
	if err != nil {
		return children, err
	}

	err = json.Unmarshal(resp.Body(), &apiResponse)
	if err != nil {
		return children, err
	}

	for _, item := range apiResponse.Data.Child {
		user := domain.User{
			ID:                  item.ID,
			ParentID:            item.ParentID,
			Image:               item.Image,
			UnitBisnis:          item.UnitBisnis,
			Name:                item.Name,
			Email:               item.Email,
			Phone:               item.Phone,
			Roles:               item.Roles,
			StatusKaryawan:      item.StatusKaryawan,
			Nik:                 item.Nik,
			Regional:            item.Regional,
			Location:            item.Location,
			Status:              item.Status,
			Salescode:           item.Salescode,
			JoinDate:            item.JoinDate,
			TglLahir:            item.TglLahir,
			Directorate:         item.Directorate,
			Department:          item.Department,
			JenisKelamin:        item.JenisKelamin,
			Privilege:           make([]string, 0),
			AdditionalPrivilege: make([]string, 0),
		}
		children = append(children, user)

	}

	return children, err
}
