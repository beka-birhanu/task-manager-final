package usercontroller

import (
	"net/http"

	"github.com/beka-birhanu/controllers/user/dto"
	usersvc "github.com/beka-birhanu/service/user"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userSvc *usersvc.Service
}

func New(userService *usersvc.Service) *UserController {
	return &UserController{userSvc: userService}
}

func (c *UserController) Register(route *gin.RouterGroup) {
	user := route.Group("/users")
	{
		user.POST("/register", c.addUser)
		user.POST("/login", c.login)
		user.PATCH("/:username/promot", c.promot)
	}
}

func (c *UserController) addUser(ctx *gin.Context) {
	var request dto.AuthRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.userSvc.Register(&usersvc.AuthCommand{
		Username: request.Username,
		Password: request.Password,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	respone := dto.NewAuthResponse(result)
	ctx.SetCookie("accessToken", result.Token, 24*60, "/", ctx.Request.Host, true, true)
	ctx.JSON(http.StatusOK, respone)
}

func (c *UserController) login(ctx *gin.Context) {
	var request dto.AuthRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.userSvc.SignIn(&usersvc.AuthCommand{
		Username: request.Username,
		Password: request.Password,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	respone := dto.NewAuthResponse(result)
	ctx.SetCookie("accessToken", result.Token, 24*60, "/", ctx.Request.Host, true, true)
	ctx.JSON(http.StatusOK, respone)
}

func (c *UserController) promot(ctx *gin.Context) {
	username, ok := ctx.Params.Get("username")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username missing"})
		return
	}

	if err := c.userSvc.Promote(username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}
