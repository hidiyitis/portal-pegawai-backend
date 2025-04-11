package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/core/usecase"
	"net/http"
)

type AttendanceHandler struct {
	attendanceUsecase usecase.AttendanceUsecase
}

func NewAttendanceHandler(attendanceUsecase usecase.AttendanceUsecase) *AttendanceHandler {
	return &AttendanceHandler{attendanceUsecase: attendanceUsecase}
}

func (h *AttendanceHandler) CreateAttendance(c *gin.Context) {
	var payload domain.InputAttendance
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	result, err := h.attendanceUsecase.CreateAttendance(payload, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated,
		gin.H{
			"data":    result,
			"message": "attendance created successfully",
		})
}
