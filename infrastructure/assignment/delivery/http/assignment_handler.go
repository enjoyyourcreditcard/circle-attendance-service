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
	userId 		:= c.Query("user_id")
	parentId 	:= c.Query("parent_id")
	startAt 	:= c.Query("start_at")
	endAt	 	:= c.Query("end_at")
	status 		:= c.Query("status")

	data, err := ash.AssignmentUsecase.GetAssignments(c.Context(), userId, parentId, startAt, endAt, status)
	if err != nil { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil}) }

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: data, Message: "SUCCESS"})
}

func (ash *AssignmentHandler) PostAssignment(c *fiber.Ctx) error {
	var assignment 		domain.Assignment
	var attachmentName 	string
	err := c.BodyParser(&assignment)
	if err != nil { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil}) }

	assignmentLocationId 	:= c.FormValue("assignment_location_id")
	startAt 				:= c.FormValue("start_at")
	endAt 					:= c.FormValue("end_at")
	parentId 				:= c.FormValue("parent_id")
	userId 					:= c.FormValue("user_id")
	attachment, err			:= c.FormFile("attachment")
	if err != nil {
		attachmentName = ""
	} else {
		fileName 		:= attachment.Filename
		extension 		:= filepath.Ext(fileName)
		hash 			:= md5.Sum([]byte(fileName))
		attachmentName	= fmt.Sprintf("%x", hash) + extension
		err 			:= c.SaveFile(attachment, "././././public/file/assignment/" + attachmentName)
		if err != nil { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil}) }
	}
	
	convertedAssignmentLocationId, err 	:= strconv.Atoi(assignmentLocationId)
	if err != nil { convertedAssignmentLocationId = 0 }
	
	convertedParentId, err := strconv.Atoi(parentId)
	if err != nil { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil}) }
	
	assignment.Attachment			= attachmentName 
	assignment.AssignmentLocationID = convertedAssignmentLocationId
	assignment.StartAt 				= startAt
	assignment.EndAt 				= endAt
	assignment.ParentID 			= convertedParentId
	assignment.UserID 				= userId
	
	err = ash.AssignmentUsecase.PostAssignment(c.Context(), &assignment)
	if err != nil { return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{Status: 500, Message: err.Error(), Data: nil}) }

	return c.JSON(domain.WebResponse{Status: http.StatusOK, Data: nil, Message: "SUCCESS"})
}