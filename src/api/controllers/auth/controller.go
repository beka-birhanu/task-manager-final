package authcontroller

import (
	"net/http"

	"github.com/beka-birhanu/task_manager_final/src/api/controllers/auth/dto"
	basecontroller "github.com/beka-birhanu/task_manager_final/src/api/controllers/base"
	errapi "github.com/beka-birhanu/task_manager_final/src/api/errors"
	icmd "github.com/beka-birhanu/task_manager_final/src/app/common/cqrs/command"
	iquery "github.com/beka-birhanu/task_manager_final/src/app/common/cqrs/query"
	registercmd "github.com/beka-birhanu/task_manager_final/src/app/user/auth/command"
	authresult "github.com/beka-birhanu/task_manager_final/src/app/user/auth/common"
	loginqry "github.com/beka-birhanu/task_manager_final/src/app/user/auth/query"
	errdmn "github.com/beka-birhanu/task_manager_final/src/domain/errors"
	"github.com/gin-gonic/gin"
)

// Controller handles HTTP requests related to authentication.
type Controller struct {
	basecontroller.BaseHandler
	registerHandler icmd.IHandler[*registercmd.Command, *authresult.Result]
	loginHandler    iquery.IHandler[*loginqry.Query, *authresult.Result]
}

// Config holds the configuration for the Controller.
type Config struct {
	RegisterHandler icmd.IHandler[*registercmd.Command, *authresult.Result]
	LoginHandler    iquery.IHandler[*loginqry.Query, *authresult.Result]
}

// New creates a new AuthController with the given CQRS handlers.
func New(config Config) *Controller {
	return &Controller{
		registerHandler: config.RegisterHandler,
		loginHandler:    config.LoginHandler,
	}
}

// RegisterPublic registers public routes.
func (c *Controller) RegisterPublic(route *gin.RouterGroup) {
	auth := route.Group("/auth")
	{
		auth.POST("/register", c.registerUser)
		auth.POST("/login", c.login)
	}
}

// RegisterProtected registers protected routes.
func (c *Controller) RegisterProtected(route *gin.RouterGroup) {
	auth := route.Group("/auth")
	{
		auth.POST("/logOut", c.logOut)
	}
}

// RegisterPrivileged registers privileged routes.
func (c *Controller) RegisterPrivileged(route *gin.RouterGroup) {}

// registerUser handles user registration.
func (c *Controller) registerUser(ctx *gin.Context) {
	var request dto.AuthRequest

	if err := ctx.ShouldBind(&request); err != nil {
		c.Problem(ctx, errapi.FromErrDMN(err.(*errdmn.Error)))
		return
	}

	result, err := c.registerHandler.Handle(registercmd.NewCommand(request.Username, request.Password))
	if err != nil {
		c.Problem(ctx, errapi.FromErrDMN(err.(*errdmn.Error)))
		return
	}

	response := dto.NewAuthResponse(result)
	c.RespondWithCookies(ctx, http.StatusOK, response, []*http.Cookie{
		{
			Name:     "accessToken",
			Value:    result.Token,
			Path:     "/",
			Domain:   ctx.Request.Host,
			MaxAge:   24 * 60 * 60,
			HttpOnly: true,
			Secure:   true,
		},
	})
}

// login handles user login.
func (c *Controller) login(ctx *gin.Context) {
	var request dto.AuthRequest

	if err := ctx.ShouldBind(&request); err != nil {
		c.Problem(ctx, errapi.NewBadRequest(err.Error()))
		return
	}

	result, err := c.loginHandler.Handle(loginqry.NewQuery(request.Username, request.Password))
	if err != nil {
		c.Problem(ctx, errapi.FromErrDMN(err.(*errdmn.Error)))
		return
	}

	response := dto.NewAuthResponse(result)
	c.RespondWithCookies(ctx, http.StatusOK, response, []*http.Cookie{
		{
			Name:     "accessToken",
			Value:    result.Token,
			Path:     "/",
			Domain:   ctx.Request.Host,
			MaxAge:   24 * 60 * 60,
			HttpOnly: true,
			Secure:   true,
		},
	})
}

// logOut handles user logout.
func (c *Controller) logOut(ctx *gin.Context) {
	c.RespondWithCookies(ctx, http.StatusNoContent, nil, []*http.Cookie{
		{
			Name:     "accessToken",
			Value:    "",
			Path:     "/",
			Domain:   ctx.Request.Host,
			MaxAge:   -1, // Delete the cookie
			HttpOnly: true,
			Secure:   true,
		},
	})
}
