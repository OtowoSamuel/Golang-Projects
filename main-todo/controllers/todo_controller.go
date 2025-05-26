package controllers

import (
	"fmt"
	"main-todo/models"
	"main-todo/utils"

	"github.com/gin-gonic/gin"
)

func CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		utils.HandleError(c, 400, "Invalid JSON data")
		return
	}

	ID := c.MustGet("ID").(uint)
	todo.ID = ID

	DB := utils.GetDB()
	if err := DB.Create(&todo).Error; err != nil {
		utils.HandleError(c, 500, "Failed to create todo")
		return
	}

	c.JSON(200, todo)
}

func GetTodos(c *gin.Context) {
	var todos []models.Todo

	// Retrieve all Todos from the database
	DB := utils.GetDB()
	DB.Find(&todos)

	c.JSON(200, todos)
}

func GetTodoByID(c *gin.Context) {
	var todo models.Todo
	todoID := c.Param("id")

	// Retrieve the Todo from the database
	DB := utils.GetDB()
	result := DB.First(&todo, todoID)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(200, todo)
}

func UpdateTodoByID(c *gin.Context) {
	var todo models.Todo
	todoID := c.Param("id")

	// Retrieve the Todo from the database
	DB := utils.GetDB()
	result := DB.First(&todo, todoID)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}

	var updatedTodo models.Todo
	if err := c.ShouldBindJSON(&updatedTodo); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON data"})
		return
	}

	// Update the Todo in the database
	todo.Title = updatedTodo.Title
	todo.Description = updatedTodo.Description
	DB.Save(&todo)

	c.JSON(200, todo)
}

func DeleteTodoByID(c *gin.Context) {
	var todo models.Todo
	todoID := c.Param("id")

	// Retrieve the Todo from the database
	DB := utils.GetDB()
	result := DB.First(&todo, todoID)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}

	// Delete the Todo from the database
	DB.Delete(&todo)

	c.JSON(200, gin.H{"message": fmt.Sprintf("Todo with ID %s deleted", todoID)})
}
