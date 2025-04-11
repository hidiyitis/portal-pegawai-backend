package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hidiyitis/portal-pegawai/internal/core/service"
	"github.com/hidiyitis/portal-pegawai/internal/core/usecase"
	"github.com/hidiyitis/portal-pegawai/internal/handler"
	"github.com/hidiyitis/portal-pegawai/internal/infrastructure/database"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
)

func StartServer() {
	db := database.NewDB()

	// Repo
	userRepo := repository.NewUserRepository(db)
	agendaRepo := repository.NewAgendaRepository(db)
	leaveRequestRepo := repository.NewLeaveRequestRepository(db)
	// Service
	userService := service.NewUserService(userRepo)
	agendaService := service.NewAgendaService(agendaRepo, userRepo)
	leaveRequestService := service.NewLeaveRequestService(leaveRequestRepo)
	// Usecase
	userUsecase := usecase.NewUserUsecase(userRepo, userService)
	agendaUsecase := usecase.NewAgendaUsecase(agendaRepo, agendaService)
	leaveRequestUsecase := usecase.NewLeaveRequestUsecase(leaveRequestRepo, leaveRequestService)
	// Handler
	userHandler := handler.NewUserHandler(userUsecase)
	agendaHandler := handler.NewAgendaHandler(agendaUsecase)
	leaveRequestHandler := handler.NewLeaveRequestHandler(leaveRequestUsecase)

	PORT := os.Getenv("PORT")
	r := gin.Default()

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server OK"})
	})

	v1 := r.Group("/api/v1")

	v1.POST("/users", userHandler.CreateUser)
	v1.POST("/auth/login", userHandler.LoginUser)
	v1.GET("/users/:nip", userHandler.GetUserByNIP)

	v1.POST("/agendas", agendaHandler.CreateAgenda)
	v1.GET("/agendas/:id", agendaHandler.GetAgendaByID)
	v1.PUT("/agendas/:id", agendaHandler.UpdateAgenda)
	v1.GET("/agendas", agendaHandler.GetAgendaByDate)
	v1.DELETE("/agendas/:id", agendaHandler.DeleteAgendaByID)

	v1.POST("/leave-request", leaveRequestHandler.CreateLeaveRequest)
	r.Run(fmt.Sprintf(":%v", PORT))
}
