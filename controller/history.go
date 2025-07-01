package controller

import (
	model "mymodule/model"
	"time"

	"github.com/sirupsen/logrus"
)

type HistoryController struct {
	Histories []model.TaskHistory
	logger logrus.Logger
}

func NewHistoryController(log *logrus.Logger) *HistoryController{
	return &HistoryController{
		Histories: []model.TaskHistory{},
		logger: *log,
	}
}

func (h *HistoryController) Record(taskID int, action string) {
	history := model.TaskHistory{
		HistoryID:        len(h.Histories) + 1,
		TaskID:    taskID,
		Action:    action,
		Timestamp: time.Now(),
	}
	h.Histories = append(h.Histories, history)

	h.logger.WithFields(logrus.Fields{
		"history_id": history.HistoryID,
		"task_id":    history.TaskID,
		"action":     action,
		"timestamp":  history.Timestamp,
	}).Info("History recorded")
}
func (h *HistoryController) GetAll() []model.TaskHistory{
	h.logger.Info("Listing all History")
	return h.Histories
}
