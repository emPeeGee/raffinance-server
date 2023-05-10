package auth

import (
	"net/http"

	"github.com/emPeeGee/raffinance/pkg/errorutil"
	"github.com/emPeeGee/raffinance/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func RegisterHandlers(authRg, apiRg *gin.RouterGroup, service Service, validate *validator.Validate, logger log.Logger) {
	h := handler{service, logger, validate}

	auth := authRg.Group("")
	{
		auth.POST("/signUp", h.signUp)
		auth.POST("/signIn", h.signIn)
	}

	api := apiRg.Group("")
	{
		api.GET("/user", h.getUser)
	}
}

type handler struct {
	service  Service
	logger   log.Logger
	validate *validator.Validate
}

func (h *handler) signUp(c *gin.Context) {
	var input createUserDTO

	if err := c.BindJSON(&input); err != nil {
		errorutil.BadRequest(c, "your request looks incorrect", err.Error())
		return
	}

	h.logger.Debug(input)

	if err := h.validate.Struct(input); err != nil {
		errorutil.BadRequest(c, "your request did not pass validation", err.Error())
		return
	}

	err := h.service.createUser(input)
	if err != nil {
		errorutil.InternalServer(c, "something went wrong, we are working", err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}

func (h *handler) signIn(c *gin.Context) {
	var credentials credentialsDTO

	if err := c.BindJSON(&credentials); err != nil {
		errorutil.BadRequest(c, "incorrect body", err.Error())
		return
	}

	if err := h.validate.Struct(credentials); err != nil {
		errorutil.BadRequest(c, "incorrect body", err.Error())
		return
	}

	token, err := h.service.generateToken(credentials)
	if err != nil {
		errorutil.InternalServer(c, "something went wrong, we are working", err.Error())
		return
	}

	h.logger.Debug(token)

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *handler) getUser(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		errorutil.Unauthorized(c, err.Error(), "")
		return
	}

	if userId == nil {
		errorutil.Unauthorized(c, "you are not authorized", "")
		return
	}

	user, err := h.service.getUserById(*userId)
	if err != nil {
		errorutil.InternalServer(c, "something went wrong, we are working", err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}
