package http

import (
	"circle/domain"
	"crypto/md5"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AssignmentHandler struct {
	AssignmentUsecase domain.AssignmentUsecase
}

func NewAssignmentHandler(app *fiber.App, asu domain.AssignmentUsecase) {
	handler := &AssignmentHandler{AssignmentUsecase: asu}

	app.Get("/assignments", handler.GetAssignments)
	app.Post("/assignment", handler.PostAssignment)

}

func (ash *AssignmentHandler) GetAssignments(c *fiber.Ctx) error {
	userId := c.Query("user_id")
	parentId := c.Query("parent_id")
	startAt := c.Query("start_at")
	endAt := c.Query("end_at")
	status := c.Query("status")
	resp, err := ash.AssignmentUsecase.GetAssignments(c.Context(), userId, parentId, startAt, endAt, status)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	if len(resp) == 0 {
		return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: make([]string, 0), Message: "SUCCESS"})
	}

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: resp, Message: "SUCCESS"})
}

func (ash *AssignmentHandler) PostAssignment(c *fiber.Ctx) error {
	var assignment domain.Assignment
	var attachmentName string
	err := c.BodyParser(&assignment)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	LocationID := c.FormValue("location_id")
	startAt := c.FormValue("start_at")
	endAt := c.FormValue("end_at")
	parentId := c.FormValue("parent_id")
	userId := c.FormValue("user_id")
	statusUrgency := c.FormValue("status_urgency")
	locationStatus := c.FormValue("location_status")
	attachment, err := c.FormFile("attachment")
	if err != nil {
		attachmentName = ""
	} else {
		fileName := attachment.Filename
		extension := filepath.Ext(fileName)
		hash := md5.Sum([]byte(fileName))
		attachmentName = fmt.Sprintf("%x", hash) + extension
		err := c.SaveFile(attachment, "././././public/file/assignment/"+attachmentName)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
		}
	}

	convertedParentID, err := strconv.Atoi(parentId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	convertedUserID, err := strconv.Atoi(userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	assignment.Attachment = attachmentName
	assignment.LocationID = LocationID
	assignment.StartAt = startAt
	assignment.EndAt = endAt
	assignment.ParentID = convertedParentID
	assignment.UserID = convertedUserID
	assignment.StatusUrgency = statusUrgency
	assignment.LocationStatus = locationStatus

	err = ash.AssignmentUsecase.PostAssignment(c.Context(), &assignment)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil})
	}

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: assignment, Message: "SUCCESS"})
}
