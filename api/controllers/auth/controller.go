package authcontroller

import (
	"net/http"

	"github.com/beka-birhanu/api/controllers/auth/dto"
	icmd "github.com/beka-birhanu/app/common/cqrs/command"
	iquery "github.com/beka-birhanu/app/common/cqrs/query"
	registercmd "github.com/beka-birhanu/app/user/auth/command"
	authresult "github.com/beka-birhanu/app/user/auth/common"
	loginqry "github.com/beka-birhanu/app/user/auth/query"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	registerHandler icmd.IHandler[*registercmd.Command, *authresult.Result]
	loginHandler    iquery.IHandler[*loginqry.Query, *authresult.Result]
}

type Config struct {
	RegisterHandler icmd.IHandler[*registercmd.Command, *authresult.Result]
	LoginHandler    iquery.IHandler[*loginqry.Query, *authresult.Result]
}

func New(config Config) *Controller {
	return &Controller{
		registerHandler: config.RegisterHandler,
		loginHandler:    config.LoginHandler,
	}
}

func (c *Controller) RegisterPublic(route *gin.RouterGroup) {
	auth := route.Group("/auth")
	{
		auth.POST("/register", c.registerUser)
		auth.POST("/login", c.login)
	}
}

func (c *Controller) RegisterProtected(route *gin.RouterGroup) {
	auth := route.Group("/auth")
	{
		auth.POST("/logOut", c.logOut)
	}

}

func (c *Controller) RegisterPrivilaged(route *gin.RouterGroup) {}

func (c *Controller) registerUser(ctx *gin.Context) {
	var request dto.AuthRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.registerHandler.Handle(registercmd.NewCommand(request.Username, request.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	respone := dto.NewAuthResponse(result)
	ctx.SetCookie("accessToken", result.Token, 24*60, "/", ctx.Request.Host, true, true)
	ctx.JSON(http.StatusOK, respone)
}

func (c *Controller) login(ctx *gin.Context) {
	var request dto.AuthRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.loginHandler.Handle(loginqry.NewQuery(request.Username, request.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	respone := dto.NewAuthResponse(result)
	ctx.SetCookie("accessToken", result.Token, 24*60, "/", ctx.Request.Host, true, true)
	ctx.JSON(http.StatusOK, respone)
}

func (c *Controller) logOut(ctx *gin.Context) {
	ctx.SetCookie("accessToken", "", 0, "/", ctx.Request.Host, true, true)
	ctx.Status(http.StatusOK)
}
