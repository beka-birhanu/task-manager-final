package common

import "github.com/gin-gonic/gin"

type IController interface {
	Register(route gin.RouterGroup)
}
