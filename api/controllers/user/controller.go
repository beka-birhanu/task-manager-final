package usercontroller

import (
	"net/http"

	basecontroller "github.com/beka-birhanu/task_manager_final/api/controllers/base"
	errapi "github.com/beka-birhanu/task_manager_final/api/errors"
	icmd "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command"
	promotcmd "github.com/beka-birhanu/task_manager_final/app/user/admin_status/command"
	errdmn "github.com/beka-birhanu/task_manager_final/domain/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Controller handles HTTP requests related to users.
type Controller struct {
	basecontroller.BaseHandler
	promotHandler icmd.IHandler[*promotcmd.Command, bool]
}

// Config holds the configuration for the Controller.
type Config struct {
	PromotHandler icmd.IHandler[*promotcmd.Command, bool]
}

// New creates a new UserController with the given CQRS handler.
func New(config Config) *Controller {
	return &Controller{
		promotHandler: config.PromotHandler,
	}
}

// RegisterPublic registers public routes.
func (c *Controller) RegisterPublic(route *gin.RouterGroup) {}

// RegisterProtected registers protected routes.
func (c *Controller) RegisterProtected(route *gin.RouterGroup) {}

// RegisterPrivileged registers privileged routes.
func (c *Controller) RegisterPrivileged(route *gin.RouterGroup) {
	user := route.Group("/users")
	{
		user.PATCH("/:username/promot", c.promot)
	}
}

// promot handles the promotion of a user.
func (c *Controller) promot(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		c.Problem(ctx, errapi.NewBadRequest("username missing"))
		return
	}

	claims, exists := ctx.Get("userClaims")
	if !exists {
		c.Problem(ctx, errapi.NewAuthentication("Claims not found"))
		return
	}

	// Type assertion to jwt.MapClaims
	jwtClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		c.Problem(ctx, errapi.NewAuthentication("Invalid claims"))
		return
	}

	// Extract and parse the user_id claim as a UUID
	userIDStr, ok := jwtClaims["user_id"].(string)
	if !ok {
		c.Problem(ctx, errapi.NewAuthentication("Invalid user_id claim"))
		return
	}

	promoterId, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Problem(ctx, errapi.NewBadRequest("Invalid user_id format"))
		return
	}

	_, err = c.promotHandler.Handle(promotcmd.NewCommand(username, promoterId))
	if err != nil {
		c.Problem(ctx, errapi.FromErrDMN(err.(*errdmn.Error)))
		return
	}

	c.Respond(ctx, http.StatusOK, nil)
}
