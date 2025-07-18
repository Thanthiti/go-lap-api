package main

import (
	
	"mymodule/Simple-Task-API/logger"
	"mymodule/Simple-Task-API/controller"
	"github.com/gin-gonic/gin"
)


func main() {
	
	log := logger.NewLogger()
	r := gin.Default()
	HistoryCtrl := controller.NewHistoryController(log)
	TaskCtrl := controller.NewTaskController(log,HistoryCtrl)

	// Task CRUD
	r.POST("/tasks/",TaskCtrl.AddTask)				// done
	r.POST("/tasks/bulk",TaskCtrl.AddTaskBulk)		// done

	r.GET("/tasks/",TaskCtrl.ShowTasks)				// done
	r.GET("/tasks/:id",TaskCtrl.GetTask)			// done
	
	r.PUT("tasks/:id",TaskCtrl.UpdateTask)			// done
	
	r.DELETE("tasks/:id",TaskCtrl.DeletedTask)		// done
	r.DELETE("tasks/bulk",TaskCtrl.DeletedTaskBulk)	// done
	
	// History
	r.GET("/histories/",TaskCtrl.ShowHistorise)
	r.Run(":8080")
}