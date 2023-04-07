package seeder

import (
	"circle/domain"
	"gorm.io/gorm"
)

func CreateAttendance(db *gorm.DB,
		UserId               string,
		Worktype             string,
		Type                 string,
		AbsenStatus          string,
		Reason               string,
		StartAt              string,
		LocationStartID      int,
		LatLongPositionStart string,
		StatusStart          string,
		Position             string,
	) error {
		return db.Create(&domain.Attendance{
			UserId               : UserId,
			Worktype             : Worktype,
			Type                 : Type,
			AbsenStatus          : AbsenStatus,
			Reason               : Reason,
			StartAt              : StartAt,
			LocationStartID      : LocationStartID,
			LatLongPositionStart : LatLongPositionStart,
			StatusStart          : StatusStart,
			Position             : Position,
		}).Error
}