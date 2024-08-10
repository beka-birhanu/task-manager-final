package usercontroller

import (
	"net/http"

	icmd "github.com/beka-birhanu/app/common/cqrs/command"
	promotcmd "github.com/beka-birhanu/app/user/admin_status/command"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	promotHandler icmd.IHandler[*promotcmd.Command, bool]
}

type Config struct {
	PromotHandler icmd.IHandler[*promotcmd.Command, bool]
}

func New(config Config) *Controller {
	return &Controller{
		promotHandler: config.PromotHandler,
	}
}

func (c *Controller) RegisterPublic(route *gin.RouterGroup) {}

func (c *Controller) RegisterProtected(route *gin.RouterGroup) {}

func (c *Controller) RegisterPrivileged(route *gin.RouterGroup) {
	user := route.Group("/users")
	{
		user.PATCH("/:username/promot", c.promot)
	}
}

func (c *Controller) promot(ctx *gin.Context) {
	username, ok := ctx.Params.Get("username")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username missing"})
		return
	}

	_, err := c.promotHandler.Handle(promotcmd.NewCommand(username, uuid.New()))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}
