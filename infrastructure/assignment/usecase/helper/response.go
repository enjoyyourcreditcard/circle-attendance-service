package helper

import (
	"circle/domain"
	"encoding/json"
	"fmt"
)

func ToAssignmentResponses(assignment []domain.Assignment) (assignmentResp []domain.AssignmentResp) {
	var assignmentResponses []domain.AssignmentResp

	for _, item := range assignment {

		var locationIDArr []int
		err := json.Unmarshal([]byte(item.LocationID), &locationIDArr)
		if err != nil {
			fmt.Println(err)
		}

		assignmentResp := domain.AssignmentResp{
			ID:             item.ID,
			Title:          item.Title,
			Description:    item.Description,
			Type:           item.Type,
			Attachment:     item.Attachment,
			LocationID:     locationIDArr,
			StartAt:        item.StartAt,
			EndAt:          item.EndAt,
			ParentID:       item.ParentID,
			UserID:         item.UserID,
			StatusUrgency:  item.StatusUrgency,
			LocationStatus: item.LocationStatus,
			Status:         item.Status,
		}

		assignmentResponses = append(assignmentResponses, assignmentResp)

	}

	return assignmentResponses
}
