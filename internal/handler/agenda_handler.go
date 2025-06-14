package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/core/usecase"
	"gorm.io/gorm"
)

type AgendaHandler struct {
	agendaUsecase usecase.AgendaUsecase
}

func NewAgendaHandler(agendaUsecase usecase.AgendaUsecase) *AgendaHandler {
	return &AgendaHandler{agendaUsecase: agendaUsecase}
}

func (h *AgendaHandler) CreateAgenda(c *gin.Context) {
	var agenda *domain.InputAgendaRequest
	if err := c.ShouldBindJSON(&agenda); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
	}

	agenda.CreatedBy = user.(domain.User).NIP

	result, err := h.agendaUsecase.CreateAgenda(agenda)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated,
		gin.H{
			"data":    result,
			"message": "agenda created successfully",
		})
}

func (h *AgendaHandler) GetAgendaByID(c *gin.Context) {
	id := c.Param("id")
	parseUint, err := strconv.ParseUint(id, 10, 32)
	result, err := h.agendaUsecase.GetAgendaByID(uint(parseUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "agenda found successfully",
	})
}

func (h *AgendaHandler) UpdateAgenda(c *gin.Context) {
	id := c.Param("id")
	parseUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	var agenda *domain.InputAgendaRequest
	if err := c.ShouldBindJSON(&agenda); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	result, err := h.agendaUsecase.UpdateAgenda(uint(parseUint), agenda)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "agenda not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "agenda updated successfully",
	})
}

func (h *AgendaHandler) DeleteAgendaByID(c *gin.Context) {
	id := c.Param("id")
	parseUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	if err = h.agendaUsecase.DeleteAgenda(uint(parseUint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "agenda deleted successfully",
	})
}

func (h *AgendaHandler) GetAgendaByDate(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "date is required"})
		return
	}
	user, ok := c.Get("user")
	nip := user.(domain.User).NIP
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
	}
	parseDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	result, _ := h.agendaUsecase.GetAgendaByDate(nip, parseDate)
	c.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "agenda found successfully",
	})
}

func (h *AgendaHandler) GetAllAgendas(c *gin.Context) {
	result, err := h.agendaUsecase.GetAllAgendas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "all agendas retrieved successfully",
	})
}
