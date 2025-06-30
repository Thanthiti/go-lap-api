package main

import (
	
	"mymodule/logger"
	"mymodule/controller"
	"github.com/gin-gonic/gin"
)


func main() {
	log := logger.NewLogger()
	r := gin.Default()
	ctrl := controller.NewTaskController(log)
	
	r.POST("/tasks/",ctrl.AddTask)				// done
	r.POST("/tasks/bulk",ctrl.AddTaskBulk)		// done
	
	r.GET("/tasks/",ctrl.ShowTasks)				// done
	r.GET("/tasks/:id",ctrl.GetTask)			// done
	
	r.PUT("tasks/:id",ctrl.UpdateTask)			// done
	
	r.DELETE("tasks/:id",ctrl.DeletedTask)		// done
	r.DELETE("tasks/bulk",ctrl.DeletedTaskBulk)	// in progress
	r.Run(":8080")
}