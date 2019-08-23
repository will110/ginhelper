package project

var baseControllerTemp = `
package controllers

import (
	"fmt"
	"github.com/ginhelper/ctl"
)

type BaseController struct {
	ctl.Controller
}

//子类初始化
type InitialiseInterface interface {
	Initialise()
}

func (m *BaseController) Prepare() {
	if app, ok := m.AppController.(InitialiseInterface); ok {
		app.Initialise()
	}
}

func (m *BaseController) Finish() {
	fmt.Println("Finish")
}
`

var userControllerTemp = `
package controllers

import "fmt"

type UserController struct {
	BaseController
}

func (m *UserController) Initialise()  {
	fmt.Println("Initialise")
}

func (m *UserController) GetList() {
	fmt.Println("GetList")
	m.C.String(200, "GetList")
}
`

var routerTemp = `
package routers

import (
	"github.com/ginhelper/routerhelper"
	"{{controllers}}/controllers"
	"{{utils}}/pkg/utils"
)

func InitRouter() {
	r := utils.NewGinDefault()
	r.GET("/user/get-list", routerhelper.BindRouter(new(controllers.UserController), "GetList"))

	r1 := r.Group("v1/user")
	{
		r1.GET("/get-list", routerhelper.BindRouter(new(controllers.UserController), "GetList"))
	}
}
`

var engineTemp = `
package utils

import "github.com/gin-gonic/gin"

var R *gin.Engine

func NewGinDefault() *gin.Engine {
	r := gin.Default()
	R = r

	return r
}
`

var mainTemp = `
package main

import (
	"log"
	"{{pkg}}/pkg/utils"
	"{{routers}}/routers"
)

func main() {
	routers.InitRouter()
	log.Fatal(utils.R.Run(":8099"))
}

`

var gitignoreTmep = `runtime/log/
debug
.idea
`