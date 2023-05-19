package api

import (
	"circle/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type apimysqlAttendanceRepository struct {
	baseUrl string
}

func NewApimysqlAttendanceRepository(baseUrl string) domain.AttendanceAPIRepository {
	return &apimysqlAttendanceRepository{baseUrl: baseUrl}
}

type ApiResponse struct {
	Status  int                      `json:"status"`
	Message string                   `json:"message"`
	Data    []domain.UserAPIResponse `json:"data"`
}

func (a apimysqlAttendanceRepository) Find(ctx context.Context, name string, nik string, unit_bisnis string, status_karyawan string, regional string) ([]domain.UserAPIResponse, error) {

	var apiResponse ApiResponse

	client := resty.New()
	var query string
	if name != "" {
		query = "&name=" + name
	}
	if nik != "" {
		query = query + "&nik=" + nik
	}
	if unit_bisnis != "" {
		query = query + "&unit_bisnis=" + unit_bisnis
	}
	if status_karyawan != "" {
		query = query + "&status_karyawan=" + status_karyawan
	}
	if regional != "" {
		query = query + "&regional=" + regional
	}

	resp, err := client.R().
		Get(a.baseUrl + "/user/api/search/user?init=init" + query)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(resp.Body(), &apiResponse)

	return apiResponse.Data, err

}
