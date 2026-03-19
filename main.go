package main

import (
	//"context"
	"log"
	"strconv"
	"w2work3/internal/config"
	"w2work3/internal/handler"
	"w2work3/internal/middleware"
	"w2work3/internal/repository"
	"w2work3/internal/service"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/logger/accesslog"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	db, err := repository.InitDB(cfg)
	if err != nil {
		panic(err)
	}
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)

	authSvc := service.NewAuthService(*userRepo, cfg.JWT.Secret, cfg.JWT.ExpireHour)
	todoSvc := service.NewTodoService(*todoRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	todoHandler := handler.NewTodoHandler(&todoSvc)

	h := server.Default(server.WithHostPorts(cfg.App.Host + ":" + strconv.Itoa(cfg.App.Port)))
	h.Use(accesslog.New())
	api := h.Group("/api/v1")

	UserGroup := api.Group("/user")
	{
		UserGroup.POST("/signup", authHandler.SignupUser)
		UserGroup.POST("/login", authHandler.LoginUser)
		UserGroup.POST("/delete", authHandler.DeleteUser)
	}

	TodoGroup := api.Group("/todo", middleware.JWTAuth(cfg.JWT.Secret))
	{
		TodoGroup.POST("/create", todoHandler.AddTodo)
		TodoGroup.POST("/list", todoHandler.ListTodo)
		TodoGroup.PATCH("/update", todoHandler.UpdateTodo)
		TodoGroup.PATCH("/status", todoHandler.UpdateTodosStatus)
		TodoGroup.DELETE("/delete", todoHandler.DeleteTodo)
	}

	log.Printf("server running at %s:%d", cfg.App.Host, cfg.App.Port)
	h.Spin()
}
