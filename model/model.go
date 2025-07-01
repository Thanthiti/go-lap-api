package model

import "time"

type Task struct {
	ID			int		  `json:ID`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Status      string    `json:"Status"`
	DueDate     time.Time `json:"DueDate"`
}

type TaskHistory struct {
    HistoryID     int       `json:"ID"`
    TaskID  int    `json:"TaskId"`
    Action string    `json:"Action"`
	Timestamp     time.Time `json:"Timestamp"`
}
