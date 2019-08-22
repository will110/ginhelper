package routerhelper

import (
	"github.com/gin-gonic/gin"
	"github.com/ginhelper/ctl"
	"reflect"
)

func BindRouter(ct ctl.ControllerInterface, methodName string) gin.HandlerFunc {
	reflectValue := reflect.ValueOf(ct)
	execController, ok := reflectValue.Interface().(ctl.ControllerInterface)
	if !ok {
		panic("controller is not ControllerInterface")
	}

	return func(c *gin.Context) {
		execController.Init(execController, c)
		execController.Prepare()
		method := reflectValue.MethodByName(methodName)
		method.Call(nil)
		execController.Finish()
	}
}
