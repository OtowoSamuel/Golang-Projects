package main

import (
	"main-todo/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/todos", controllers.CreateTodo)
	r.GET("/todos", controllers.GetTodos)
	r.GET("/todos/:id", controllers.GetTodoByID)
	r.PUT("/todos/:id", controllers.UpdateTodoByID)
	r.DELETE("/todos/:id", controllers.DeleteTodoByID)

	r.Run(":8080") // listen and serve on

}
