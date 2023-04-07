package seeder

import "gorm.io/gorm"

func Seed(db *gorm.DB) {
	// Create attendances
	CreateAttendance(db, "1", "wfo", "office_hour", "", "Ujian Sekolah", "07-04-2023 11:21:20", 1, "lat: -6.372651233737751, long: 106.83478630771553", "inside_area", "")
}