package controller

import (
	model "mymodule/model"
	response "mymodule/response"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BulkDeleteRequest struct{
	IDs []int `json:"ids"` 	
}
type TaskController struct {
	Tasks  []model.Task
	nextID int
	history *HistoryController
	logger logrus.Logger
}

func NewTaskController(log *logrus.Logger,history *HistoryController) *TaskController {
	return &TaskController{
		Tasks:  []model.Task{},
		nextID: 1,
		history: history,
		logger: *log,
	}
}

func (c *TaskController) AddTask(ctx *gin.Context) {
	var t model.Task
	if err := ctx.ShouldBindJSON(&t); err != nil {
		c.logger.WithError(err).Error("Invalid JSON")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("ERR_INVALID_JSON", "Invalid JSON"))
		return
	}

	// Title required
	if strings.TrimSpace(t.Title) == "" {
		c.logger.Warn("Title is required")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("ERR_TITLE_REQUIRED", "TITLE is required"))
		return
	}

	// Status  pending, in-progress, done and required
	status := map[string]bool{
		"pending":     true,
		"in-progress": true,
		"done":        true,
	}

	if strings.TrimSpace(t.Status) == "" {
		c.logger.Warn("Status is required")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("ERR_STATUS_REQUIRED", "Status is required"))
		return

	}
	if !status[t.Status] {
		c.logger.Warn("Status is wrong")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("ERR_STATUS_WRONG", "Status is wrong"))
		return
	}

	// Date Validation
	if t.DueDate.Before(time.Now()) || t.DueDate.IsZero() {
		c.logger.Warn("Date is wrong")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("ERR_DATE_WRONG", "Date is wrong"))
		return
	}
	
	// Description Validation
	if strings.TrimSpace(t.Description) == "" {
		t.Description = "-"
	}

	t.ID = c.nextID
	c.nextID++
	c.Tasks = append(c.Tasks, t)
	
	c.logger.WithFields(logrus.Fields{
		"id":          t.ID,
		"title":       t.Title,
		"status":      t.Status,
		"description": t.Description,
		"dueDate":     t.DueDate,
	}).Info("Added new Todo")
	c.history.Record(t.ID,"Add")
	ctx.JSON(http.StatusOK, t)
}

func (c *TaskController) AddTaskBulk(ctx *gin.Context) {
	var t []model.Task
	if err := ctx.ShouldBindJSON(&t); err != nil {
		c.logger.WithError(err).Error("Invalid JSON")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("ERR_INVALID_JSON", "Invalid JSON"))
		return
	}
	// Title required
	for  i := range t {

		if strings.TrimSpace(t[i].Title) == "" {
			c.logger.Warn("Title is required")
			ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("ERR_TITLE_REQUIRED", "TITLE is required"))
			return
		}
		// Status  pending, in-progress, done and required
		status := map[string]bool{
			"pending":     true,
			"in-progress": true,
			"done":        true,
		}

		if strings.TrimSpace(t[i].Status) == "" {
			c.logger.Warn("Status is required")
			ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("ERR_STATUS_REQUIRED", "Status is required"))
			return

		}
		if !status[t[i].Status] {
			c.logger.Warn("Status is wrong")
			ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("ERR_STATUS_WRONG", "Status is wrong"))
			return
		}
		// Date Validation
		if t[i].DueDate.Before(time.Now()) || t[i].DueDate.IsZero() {
			c.logger.Warn("Date is wrong")
			ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("ERR_DATE_WRONG", "Date is wrong"))
			return
		}
		// Description Validation
		if strings.TrimSpace(t[i].Description) == "" {
			t[i].Description = "-"
		}
		t[i].ID = c.nextID
		c.nextID++
		c.Tasks = append(c.Tasks, t[i])

		c.logger.WithFields(logrus.Fields{
			"id":          t[i].ID,
			"title":       t[i].Title,
			"status":      t[i].Status,
			"description": t[i].Description,
			"dueDate":     t[i].DueDate,
		}).Info("Added new Todo")
		c.history.Record(t[i].ID,"Add")
	}
	ctx.JSON(http.StatusOK, t)
}

func (c *TaskController) ShowTasks(ctx *gin.Context) {
	c.logger.Info("Listing all Tasks")
	ctx.JSON(http.StatusOK, c.Tasks)
}

func (c *TaskController) ShowHistorise(ctx *gin.Context) {
	c.logger.Info("Listing all Tasks")
	ctx.JSON(http.StatusOK, c.history.GetAll())
}

func (c *TaskController) GetTask(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id ,err := strconv.Atoi(idstr)
	if err != nil || id <= 0{
		c.logger.WithField("idParam",idstr).Error("Invalid ID")
        ctx.JSON(http.StatusBadRequest,response.NewErrorResponse("ERR_INVALID_ID", "Invalid ID"))
        return
	}
	
	found := false
	for _,task:= range c.Tasks{
		if task.ID == id{
			ctx.JSON(http.StatusOK,task)
			found = true
			c.logger.WithField("id",id).Info("Task found")
			break
		}
	}
	if !found {
		c.logger.WithField("id", id).Warn("Task not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("ERR_TASK_NOT_FOUND", "Task not found"))
		return
	}
	
	
}

func (c *TaskController) UpdateTask(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id ,err := strconv.Atoi(idstr)
	if err != nil || id <= 0{
		c.logger.WithField("idParam",idstr).Error("Invalid ID")
		ctx.JSON(http.StatusBadRequest,response.NewErrorResponse("ERR_INVALID_ID", "Invalid ID"))
		return
	}
	var updateTask model.Task
	if err := ctx.ShouldBindJSON(&updateTask); err != nil {
		c.logger.WithError(err).Error("Invalid JSON")
		ctx.JSON(http.StatusBadRequest,response.NewErrorResponse("ERR_INVALID_JSON", "Invalid JSON"))
        return
    }
	if strings.TrimSpace(updateTask.Title) == "" {
		c.logger.Warn("Task is required")
		ctx.JSON(http.StatusBadRequest,response.NewErrorResponse("ERR_TASK_REQUIRED", "Task is required"))
        return
    }
	updated := false
    for i, todo := range c.Tasks {
        if todo.ID == id {
            c.Tasks[i].Title = updateTask.Title
            updated = true
			c.logger.WithFields(logrus.Fields{
				"id":id,
				"task":updateTask.Title,
			}).Info("Todo updated")
            ctx.JSON(http.StatusOK, c.Tasks[i])
            break
        }
    }
	
    if !updated {
		c.logger.WithField("id", id).Warn("Todo not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("ERR_TODO_NOT_FOUND", "Todo not found"))
		return
	}

	
}

func (c *TaskController) DeletedTask(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.logger.WithField("idParam",idStr).Error("Invalid ID")
        ctx.JSON(http.StatusBadRequest,response.NewErrorResponse("ERR_INVALID_ID", "Invalid ID"))        
		return
	}

	found := false
	for i, data := range c.Tasks {
		if data.ID == id {
			c.Tasks = append(c.Tasks[:i], c.Tasks[i+1:]...)
			found = true

			c.logger.WithField("id",id).Info("Todo Deleted")
			break
		}
	}

	if !found {
		c.logger.WithField("id", id).Warn("Todo not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("ERR_TODO_NOT_FOUND", "Todo not found"))
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *TaskController) DeletedTaskBulk(ctx *gin.Context) {
	var req BulkDeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil{
		c.logger.WithError(err).Error("Invalid JSON")
		ctx.JSON(http.StatusBadRequest,response.NewErrorResponse("ERR_INVALID_JSON","Invalid JSON"))
		return
	}

	
	// NewTask := []model.Task{}
	 var NewTask []model.Task
	 deleteID := []int{}
	for _,task := range c.Tasks{
		found := false
		for _,id := range req.IDs {
			if id == task.ID{
				found = true
				break
			}	
		}
		if !found{
			NewTask = append(NewTask, task)
		}else{
			deleteID = append(deleteID, task.ID)	
			c.logger.WithField("id",task.ID).Info("Task deleted (bulk)")
		}

	}
	c.logger.WithFields(logrus.Fields{
		"deletedCount": len(deleteID),
		"deletedIDs":   deleteID,
		"remaining":    len(c.Tasks),
	}).Info("Bulk delete completed")
	c.Tasks = NewTask
	ctx.Status(http.StatusNoContent)
	
}
