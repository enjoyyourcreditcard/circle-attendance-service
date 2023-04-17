package domain

import (
	"context"
)

type Attendance struct {
	ID                   int    `json:"id" gorm:"primaryKey"`
	UserId               string `json:"user_id" gorm:"type:varchar(255)"`
	Absen                int    `json:"absen"`
	StartAt              string `json:"start_at" gorm:"type:varchar(255)"`
	UpdatedAt            string `json:"updated_at" gorm:"type:varchar(255)"`
	EndAt                string `json:"end_at" gorm:"type:varchar(255)"`
	Timezone             string `json:"timezone" gorm:"type:varchar(255)"`
	Position             string `json:"position" gorm:"type:varchar(1000)"`
	Worktype             string `json:"worktype" gorm:"type:varchar(255)"`
	LocationStartID      int    `json:"location_start_id"`
	LocationEndID        int    `json:"location_end_id"`
	Type                 string `json:"type" gorm:"type:varchar(255)"`
	WorkingHour          string `json:"working_hour" gorm:"type:varchar(255)"`
	Reason               string `json:"reason" gorm:"type:varchar(255)"`
	AbsenStatus          string `json:"absen_status" gorm:"type:text"`
	Notes                string `json:"notes" gorm:"type:varchar(255)"`
	LatLongPositionStart string `json:"lat_long_position_start" gorm:"type:text"`
	LatLongPositionEnd   string `json:"lat_long_position_end" gorm:"type:text"`
	StatusStart          string `json:"status_start" gorm:"type:varchar(255)"`
	StatusEnd            string `json:"status_end" gorm:"type:varchar(255)"`
	CreatedAt            string `json:"created_at" gorm:"type:varchar(255)"`
}

type EndAttendance struct {
	EndAt              	string `json:"end_at" gorm:"type:varchar(255)"`
	Reason				string `json:"reason" gorm:"type:varchar(255)"`
	LocationEndID      	int    `json:"location_end_id"`
	LatLongPositionEnd 	string `json:"lat_long_position_end" gorm:"type:text"`
	StatusEnd          	string `json:"status_end" gorm:"type:varchar(255)"`
	WorkingHour        	string `json:"working_hour" gorm:"type:varchar(255)"`
}

type DashboardAttendance struct {
	WorkingDay 				int64 `json:"working_day"`
	NonWorkingDay 			int64 `json:"non_working_day"`
	Holiday					int64 `json:"holiday"`
	TotalClockin 			int64 `json:"total_clockin"`
	TotalClockout 			int64 `json:"total_clockout"`
	TotalWfh				int64 `json:"total_wfh"`
	TotalWfo				int64 `json:"total_wfo"`
	LateIn					int64 `json:"late_in"`
	EarlyIn				int64 `json:"early_in"`
	EarlyOut				int64 `json:"early_out"`
	InsideArea				int64 `json:"inside_area"`
	OutsideArea				int64 `json:"outside_area"`
	InsideOtherArea 		int64 `json:"inside_other_area"`
	Shifting				int64 `json:"shifting"`
	OfficeHour				int64 `json:"office_hour"`
	Alpa					int64 `json:"apla"`
	Sick					int64 `json:"sick"`
	Izin					int64 `json:"izin"`
	Leave					int64 `json:"leave"`
	UneligibleWorkingHour	int64 `json:"uneligible_working_hour"`
	Penugasan				int64 `json:"penugasan"`
}

type AttendanceUsecase interface {
	GetUserLastAttendance		(context.Context, string)					(Attendance, error)
	GetUserDashboardAttendance	(context.Context, string, string, string)	(DashboardAttendance, error)
	GetUserAttendanceData		(context.Context, string, string, string)	([]Attendance, error)

	PostClockIn				(context.Context, *Attendance, string) 		(string, error)
	PostClockOut			(context.Context, *EndAttendance, string) 	(string, error)
	PostAttendanceNotes		(context.Context, string, string) 			(string, error)
}

type AttendanceRepository interface {
	CheckAbsen					(context.Context, string, string) 			(int, error)
	GetUserLastAttendance		(context.Context, string)					(Attendance, error)
	GetUserDashboardAttendance	(context.Context, string, string, string)	(DashboardAttendance, error)
	GetUserAttendanceData		(context.Context, string, string, string)	([]Attendance, error)

	CreateAbsen				(context.Context, *Attendance) 			(*Attendance, error)
	UpdateAbsen				(context.Context, *EndAttendance, int) 	error
	PostAttendanceNotes		(context.Context, string, string) 		error
}