package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/core/usecase"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.userUsecase.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated,
		gin.H{
			"data":    user,
			"message": "user created successfully",
		})
}

func (h *UserHandler) GetUserByNIP(c *gin.Context) {
	nip := c.Param("nip")
	if nip == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid ID"})
		return
	}
	u64, _ := strconv.ParseUint(nip, 10, 64)
	user, err := h.userUsecase.GetUserByNIP(uint(u64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"data":    user,
			"message": "success find user",
		})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	type LoginUser struct {
		Nip      string `json:"nip" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	loginUser := LoginUser{}
	var user domain.User
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	nip, _ := strconv.ParseUint(loginUser.Nip, 10, 64)
	user = domain.User{
		NIP:      uint(nip),
		Password: loginUser.Password,
	}
	accessToken, refreshToken, expiredAt, err := h.userUsecase.LoginUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":                    user,
			"access_token":            accessToken,
			"refresh_token":           refreshToken,
			"access_token_expired_at": expiredAt,
		},
		"message": "login successful",
	})
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, isExist := c.Get("user")
	if !isExist {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		return
	}

	userObj := user.(domain.User)
	result, err := h.userUsecase.UploadAvatar(c.Request.Context(), &userObj, fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"data":    result,
	})
}

func (h *UserHandler) UpdateUserPassword(c *gin.Context) {
	updateUserPassword := domain.UpdateUserPassword{}
	if err := c.ShouldBindJSON(&updateUserPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, isExist := c.Get("user")
	if !isExist {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		return
	}
	userObj := user.(domain.User)
	user, err := h.userUsecase.UpdateUserPassword(userObj, updateUserPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "user update password successfully",
	})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userUsecase.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    users,
		"message": "all users retrieved successfully",
	})
}

func (h *UserHandler) GetUsersExclude(c *gin.Context) {
	user, isExist := c.Get("user")
	if !isExist {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		return
	}
	userObj := user.(domain.User)
	result, err := h.userUsecase.GetUsersExclude(userObj.NIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "success get users",
	})
}
