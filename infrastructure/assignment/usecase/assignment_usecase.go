package usecase

import (
	"circle/domain"
	"context"
	"time"
)

type assignmentUsecase struct {
	assignmentRepo domain.AssignmentRepository
	contextTimeout time.Duration
}

func NewAssignmentUsecase(as domain.AssignmentRepository, timeout time.Duration) domain.AssignmentUsecase {
	return &assignmentUsecase{
		assignmentRepo: as,
		contextTimeout: timeout,
	}
}

func (as assignmentUsecase) GetAssignments(ctx context.Context, userId string, parentId string, startAt string, endAt string, status string) ([]domain.Assignment, error) {
	res, err := as.assignmentRepo.GetAssignments(ctx, userId, parentId, startAt, endAt, status)
	return res, err
}

func (as assignmentUsecase) PostAssignment(ctx context.Context, assignment *domain.Assignment) error {
	err := as.assignmentRepo.CreateAssignment(ctx, assignment)
	return err
}