package domain

import "context"

type Assignment struct {
	ID                   int    `json:"id" gorm:"primaryKey"`
	Title                string `json:"title" gorm:"type:varchar(255)"`
	Description          string `json:"description" gorm:"type:text"`
	Type                 string `json:"type" gorm:"type:varchar(255)"`
	Attachment           string `json:"attachment" gorm:"type:varchar(255)"`
	AssignmentLocationID int    `json:"assignment_location_id"`
	StartAt              string `json:"start_at" gorm:"type:varchar(255)"`
	EndAt                string `json:"end_at" gorm:"type:varchar(255)"`
	ParentID             int    `json:"parent_id"`
	UserID               string `json:"user_id"`
	Status               string `json:"status" gorm:"type:varchar(255)"`
}

type AssignmentUsecase interface {
	GetAssignments(context.Context, string, string, string, string, string) ([]Assignment, error)
	
	PostAssignment(context.Context, *Assignment) error
}

type AssignmentRepository interface {
	GetAssignments(context.Context, string, string, string, string, string) ([]Assignment, error)

	CreateAssignment(context.Context, *Assignment) error
}