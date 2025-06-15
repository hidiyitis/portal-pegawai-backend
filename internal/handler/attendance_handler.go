package handler

import (
	"context"
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
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
	}
	userObj := user.(domain.User)

	payload.NIP = userObj.NIP
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	result, err := h.attendanceUsecase.CreateAttendance(context.Background(), payload, fileHeader)
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

func (h *AttendanceHandler) GetLastAttendance(c *gin.Context) {
	user, isExist := c.Get("user")
	if !isExist {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		return
	}
	userObj := user.(domain.User)
	result, err := h.attendanceUsecase.GetLastAttendance(userObj.NIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "all agendas retrieved successfully",
	})
}
