package controller

import (
	"blog/dao"
	"blog/model"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(200, "register.html", nil)
}
func AddUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	user := model.User{
		Username: username,
		Password: password,
	}
	dao.Mgr.AddUser(&user)
}
