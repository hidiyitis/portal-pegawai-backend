package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/core/usecase"
	"net/http"
)

type LeaveRequestHandler struct {
	leaveRequestUsercase usecase.LeaveRequestUsecase
}

func NewLeaveRequestHandler(repo usecase.LeaveRequestUsecase) *LeaveRequestHandler {
	return &LeaveRequestHandler{
		leaveRequestUsercase: usecase.LeaveRequestUsecase(repo),
	}
}

func (h *LeaveRequestHandler) CreateLeaveRequest(c *gin.Context) {
	var leaveRequest *domain.InputLeaveRequest
	if err := c.Bind(&leaveRequest); err != nil {
		fmt.Println("err")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println("err")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = h.leaveRequestUsercase.CreateLeaveRequest(leaveRequest, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, leaveRequest)
}
