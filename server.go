package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title     string `json:"title"`
	Completed int    `json:"completed"`
	Description string `json:"description"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Todo{})

	e := echo.New();

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	
	e.GET("/todos", func(c echo.Context) error {
		todos := []Todo{}
		db.Find(&todos)
		return c.JSON(http.StatusOK, todos)
	})
	e.POST("/todos", func(c echo.Context) error {
		todo := new(Todo)
		if err := c.Bind(todo); err != nil {
			return err
		}
		db.Create(&todo)
		return c.JSON(http.StatusCreated, todo)
	})
	e.PUT("/todos/:id", func(c echo.Context) error {
		id := c.Param("id")
		todo := new(Todo)
		if err := c.Bind(todo); err != nil {
			return err
		}
		db.Model(&todo).Where("id = ?", id).Updates(Todo{Title: todo.Title, Completed: todo.Completed, Description: todo.Description})
		return c.JSON(http.StatusOK, todo)
	})
	e.DELETE("/todos/:id", func(c echo.Context) error {
		id := c.Param("id")
		db.Where("id = ?", id).Delete(&Todo{})
		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":4000"))
}