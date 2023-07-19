package router

import (
	"blog/controller"

	"github.com/gin-gonic/gin"
)

func Work() {
	e := gin.Default()
	e.LoadHTMLGlob("templates/*")
	e.Static("/assets", "./assets")

	e.GET("/register", controller.Index)
	e.POST("/register", controller.AddUser)
	e.Run()
}
