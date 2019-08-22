package ctl

import "github.com/gin-gonic/gin"

type ControllerInterface interface {
	Init(ct ControllerInterface, c *gin.Context)
	Prepare()
	Finish()
}
