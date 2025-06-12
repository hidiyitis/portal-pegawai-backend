package http

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hidiyitis/portal-pegawai/internal/infrastructure/storage"
	"github.com/hidiyitis/portal-pegawai/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/hidiyitis/portal-pegawai/internal/core/service"
	"github.com/hidiyitis/portal-pegawai/internal/core/usecase"
	"github.com/hidiyitis/portal-pegawai/internal/handler"
	"github.com/hidiyitis/portal-pegawai/internal/infrastructure/database"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
)

func StartServer() {

	ctx := context.Background()
	db := database.NewDB()
	gcpStorage, err := storage.NewGCPStorage(ctx, "")
	if err != nil {
		panic(err)
	}

	// Repo
	userRepo := repository.NewUserRepository(db)
	agendaRepo := repository.NewAgendaRepository(db)
	leaveRequestRepo := repository.NewLeaveRequestRepository(db)
	holidayRepo := repository.NewHolidayRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)
	// Service
	userService := service.NewUserService(userRepo, gcpStorage)
	agendaService := service.NewAgendaService(agendaRepo, userRepo)
	leaveRequestService := service.NewLeaveRequestService(leaveRequestRepo, userRepo, holidayRepo, gcpStorage)
	attendanceService := service.NewAttendanceService(attendanceRepo, gcpStorage)
	// Usecase
	userUsecase := usecase.NewUserUsecase(userRepo, userService)
	agendaUsecase := usecase.NewAgendaUsecase(agendaRepo, agendaService)
	leaveRequestUsecase := usecase.NewLeaveRequestUsecase(leaveRequestRepo, leaveRequestService)
	attendanceUsecase := usecase.NewAttendanceUsecase(attendanceRepo, attendanceService)
	// Handler
	userHandler := handler.NewUserHandler(userUsecase)
	agendaHandler := handler.NewAgendaHandler(agendaUsecase)
	leaveRequestHandler := handler.NewLeaveRequestHandler(leaveRequestUsecase)
	attendanceHandler := handler.NewAttendanceHandler(attendanceUsecase)

	PORT := os.Getenv("PORT")
	r := gin.Default()

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server OK"})
	})

	v1 := r.Group("/api/v1")

	v1.POST("/users", userHandler.CreateUser)
	v1.GET("/users", utils.AuthMiddleware(), userHandler.GetUsers)
	v1.POST("/auth/login", userHandler.LoginUser)
	v1.GET("/users/", utils.AuthMiddleware(), userHandler.GetUsersExclude)
	v1.GET("/users/:nip", utils.AuthMiddleware(), userHandler.GetUserByNIP)
	v1.PUT("/users/upload-avatar", utils.AuthMiddleware(), userHandler.UploadAvatar)
	v1.PUT("/users/update-password", utils.AuthMiddleware(), userHandler.UpdateUserPassword)

	v1.POST("/agendas", utils.AuthMiddleware(), agendaHandler.CreateAgenda)
	v1.GET("/agendas/:id", utils.AuthMiddleware(), agendaHandler.GetAgendaByID)
	v1.PUT("/agendas/:id", utils.AuthMiddleware(), agendaHandler.UpdateAgenda)
	v1.GET("/agendas", utils.AuthMiddleware(), agendaHandler.GetAgendaByDate)
	v1.DELETE("/agendas/:id", utils.AuthMiddleware(), agendaHandler.DeleteAgendaByID)
	v1.GET("/agendas/all", utils.AuthMiddleware(), agendaHandler.GetAllAgendas)

	v1.POST("/leave-request", utils.AuthMiddleware(), leaveRequestHandler.CreateLeaveRequest)
	v1.GET("/leave-request", utils.AuthMiddleware(), leaveRequestHandler.GetLeaveRequest)
	v1.GET("/dashboard-leave-request", utils.AuthMiddleware(), leaveRequestHandler.GetDashboardLeaveRequest)
	v1.PUT("/leave-request/:id", utils.AuthMiddleware(), leaveRequestHandler.UpdateLeaveRequest)

	v1.POST("/attandance", utils.AuthMiddleware(), attendanceHandler.CreateAttendance)
	r.Run(fmt.Sprintf(":%v", PORT))
}
