package domain

type Attendance struct {
	ID        				int			`json:"id" gorm:"primaryKey"`
	UserId    				string		`json:"user_id" gorm:"type:varchar(255)"`
	Absen     				int			`json:"absen"`
	StartAt   				string		`json:"start_at" gorm:"type:varchar(255)"`
	UpdatedAt 				string		`json:"updated_at" gorm:"type:varchar(255)"`
	EndAt 					string		`json:"end_at" gorm:"type:varchar(255)"`
	Position 				string		`json:"position" gorm:"type:varchar(1000)"`
	Worktype 				string		`json:"worktype" gorm:"type:varchar(255)"`
	LocationStartID 		int			`json:"location_start_id"`
	LocationEndID 			int			`json:"location_end_id"`
	Type 					string		`json:"type" gorm:"type:varchar(255)"`
	WorkingHour 			string		`json:"working_hour" gorm:"type:varchar(255)"`
	Reason 					string		`json:"reason" gorm:"type:varchar(255)"`
	AbsenStatus 			string		`json:"absen_status" gorm:"type:text"`
	Notes 					string		`json:"notes" gorm:"type:varchar(255)"`
	LatLongPositionStart 	string		`json:"lat_long_position_start" gorm:"type:text"`
	LatLongPositionEnd 		string		`json:"lat_long_position_end" gorm:"type:text"`
	StatusStart 			string		`json:"status_start" gorm:"type:varchar(255)"`
	StatusEnd 				string		`json:"status_end" gorm:"type:varchar(255)"`
	CreatedAt 				string		`json:"created_at" gorm:"type:varchar(255)"`
	Timezone 				string		`json:"timezone" gorm:"type:varchar(255)"`
}

type AttendanceUsecase interface {
}

type AttendanceRepository interface {
}
