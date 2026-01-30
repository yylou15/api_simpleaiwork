package handler

import (
	"net/http"
	"strconv"

	"api/biz/say_right/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	svc service.TemplateService
}

func NewTemplateHandler(svc service.TemplateService) *TemplateHandler {
	return &TemplateHandler{
		svc: svc,
	}
}

func (h *TemplateHandler) ListTemplates(c *gin.Context) {
	userID, ok := getSessionUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	categories, err := h.svc.ListTemplatesByCategory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

func (h *TemplateHandler) GetTemplateDetail(c *gin.Context) {
	userID, ok := getSessionUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idParam := c.Param("id")
	idValue, err := strconv.Atoi(idParam)
	if err != nil || idValue <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template id"})
		return
	}

	result, err := h.svc.GetTemplateDetail(c.Request.Context(), userID, int32(idValue))
	if err != nil {
		if err == service.ErrProRequired {
			c.JSON(http.StatusForbidden, gin.H{"error": "Pro required"})
			return
		}
		if err == service.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func getSessionUserID(c *gin.Context) (int32, bool) {
	session := sessions.Default(c)
	value := session.Get("user_id")
	if value == nil {
		return 0, false
	}
	switch v := value.(type) {
	case int32:
		return v, v > 0
	case int:
		return int32(v), v > 0
	case int64:
		return int32(v), v > 0
	case float64:
		return int32(v), v > 0
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		return int32(parsed), parsed > 0
	default:
		return 0, false
	}
}
