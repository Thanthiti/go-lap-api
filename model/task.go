package model

import "time"

type Task struct {
	ID			int		  `json:ID`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Status      string    `json:"Status"`
	DueDate     time.Time `json:"DueDate"`
}