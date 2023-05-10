package tag

import (
	"net/http"
	"strconv"

	"github.com/emPeeGee/raffinance/internal/auth"
	"github.com/emPeeGee/raffinance/pkg/errorutil"
	"github.com/emPeeGee/raffinance/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func RegisterHandlers(apiRg *gin.RouterGroup, service Service, validate *validator.Validate, logger log.Logger) {
	h := handler{service, logger, validate}

	api := apiRg.Group("/tags")
	{
		api.GET("", h.getTags)
		api.POST("", h.createTag)
		api.PUT("/:id", h.updateTag)
		api.DELETE("/:id", h.deleteTag)

	}
}

type handler struct {
	service  Service
	logger   log.Logger
	validate *validator.Validate
}

func (h *handler) createTag(c *gin.Context) {
	var input createTagDTO

	userId, err := auth.GetUserId(c)
	if err != nil || userId == nil {
		errorutil.Unauthorized(c, err.Error(), "you are not authorized")
		return
	}

	if err := c.BindJSON(&input); err != nil {
		errorutil.BadRequest(c, "your request looks incorrect", err.Error())
		return
	}

	h.logger.Debug(input)

	if err := h.validate.Struct(input); err != nil {
		errorutil.BadRequest(c, "your request did not pass validation", err.Error())
		return
	}

	createdTag, err := h.service.createTag(*userId, input)
	if err != nil {
		errorutil.InternalServer(c, "It looks like name is already used", err.Error())
		return
	}

	c.JSON(http.StatusOK, createdTag)
}

func (h *handler) updateTag(c *gin.Context) {
	var input updateTagDTO

	userId, err := auth.GetUserId(c)
	if err != nil || userId == nil {
		errorutil.Unauthorized(c, err.Error(), "you are not authorized")
		return
	}

	tagId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errorutil.BadRequest(c, "wrong tag id", err.Error())
		return
	}

	if err := c.BindJSON(&input); err != nil {
		errorutil.BadRequest(c, "your request looks incorrect", err.Error())
		return
	}

	if err := h.validate.Struct(input); err != nil {
		errorutil.BadRequest(c, "your request did not pass validation", err.Error())
		return
	}

	updatedTag, err := h.service.updateTag(*userId, uint(tagId), input)
	if err != nil {
		errorutil.BadRequest(c, "error", err.Error())
		return
	}

	c.JSON(http.StatusOK, updatedTag)
}

func (h *handler) deleteTag(c *gin.Context) {
	userId, err := auth.GetUserId(c)
	if err != nil || userId == nil {
		errorutil.Unauthorized(c, err.Error(), "you are not authorized")
		return
	}

	tagId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errorutil.BadRequest(c, err.Error(), "the id must be an integer")
		return
	}

	if err := h.service.deleteTag(*userId, uint(tagId)); err != nil {
		h.logger.Info(err.Error())
		errorutil.NotFound(c, err.Error(), "Not found")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}

func (h *handler) getTags(c *gin.Context) {
	userId, err := auth.GetUserId(c)
	if err != nil || userId == nil {
		errorutil.Unauthorized(c, err.Error(), "you are not authorized")
		return
	}

	tags, err := h.service.getTags(*userId)
	if err != nil {
		errorutil.InternalServer(c, "something went wrong, we are working", err.Error())
		return
	}

	c.JSON(http.StatusOK, tags)
}
