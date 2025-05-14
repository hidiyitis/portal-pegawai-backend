package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/core/usecase"
	"net/http"
	"strconv"
)

type LeaveRequestHandler struct {
	leaveRequestUsercase usecase.LeaveRequestUsecase
}

func NewLeaveRequestHandler(leaveRequestUsercase usecase.LeaveRequestUsecase) *LeaveRequestHandler {
	return &LeaveRequestHandler{
		leaveRequestUsercase: leaveRequestUsercase,
	}
}

func (h *LeaveRequestHandler) CreateLeaveRequest(c *gin.Context) {
	var leaveRequest *domain.InputLeaveRequest
	if err := c.Bind(&leaveRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"message": "user not found"})
	}
	userObj := user.(domain.User)

	result, err := h.leaveRequestUsercase.CreateLeaveRequest(c.Request.Context(), userObj, leaveRequest, fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *LeaveRequestHandler) UpdateLeaveRequest(c *gin.Context) {
	var leaveRequest *domain.LeaveRequest

	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"message": "user not found"})
	}
	userObj := user.(domain.User)

	if err := c.Bind(&leaveRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	id := c.Param("id")
	parseUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := h.leaveRequestUsercase.UpdateLeaveRequest(uint(parseUint), userObj, leaveRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)

}

func (h *LeaveRequestHandler) GetLeaveRequest(c *gin.Context) {
	status := c.Query("status")
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"message": "user not found"})
		return
	}
	userObj := user.(domain.User)
	result, err := h.leaveRequestUsercase.GetLeaveRequests(userObj.NIP, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result, "message": "success get leave requests"})
}

func (h *LeaveRequestHandler) GetDashboardLeaveRequest(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"message": "user not found"})
		return
	}
	userObj := user.(domain.User)
	result, err := h.leaveRequestUsercase.GetLeaveRequestDashboard(userObj.NIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result, "message": "success dashboard leave requests"})
}
