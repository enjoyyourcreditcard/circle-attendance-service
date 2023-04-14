package mysql

import (
	"circle/domain"
	"context"

	"gorm.io/gorm"
)

type mysqlAssignmentRepository struct {
	conn *gorm.DB
}

func NewMysqlAssignmentRepository(conn *gorm.DB) domain.AssignmentRepository {
	return &mysqlAssignmentRepository{conn}
}

func (ar mysqlAssignmentRepository) GetAssignments(ctx context.Context, userId string, parentId string, startAt string, endAt string, status string) ([]domain.Assignment, error) {
	var assignments []domain.Assignment
	result := ar.conn.Where("user_id LIKE ?", "%"+userId+"%").Where("parent_id LIKE ?", "%"+parentId+"%").Where("start_at LIKE ?", "%"+startAt+"%").Where("end_at LIKE ?", "%"+endAt+"%").Where("status LIKE ?", "%"+status+"%").Find(&assignments)
	return assignments, result.Error
}

func (ar mysqlAssignmentRepository) CreateAssignment(ctx context.Context, assignment *domain.Assignment) error {
	result := ar.conn.Create(assignment)
	return result.Error
}