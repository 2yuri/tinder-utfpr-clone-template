package api

import (
	"context"
	"tinderutf/api/handlers"
	"tinderutf/api/middlewares"
	"tinderutf/api/websocket"
)

func (s *Server) setupRoutes() {
	s.app.Use(middlewares.EnableCors())

	s.app.Static("/img", "./images")

	v1 := s.app.Group("api/v1")

	v1.POST("/login", handlers.Login)
	v1.GET("/me", middlewares.Auth(), handlers.Me)

	user := v1.Group("user")
	{
		user.POST("/", handlers.CreateUser)
		user.POST("/customize", middlewares.Auth(), handlers.SetCustomization)
		user.POST("/geolocation", middlewares.Auth(), handlers.SetGeolocation)
		user.POST("/image", middlewares.Auth(), handlers.SetProfileImage)
		user.GET("/:id", middlewares.Auth(), handlers.GetById)
	}

	interactions := v1.Group("interactions", middlewares.Auth())
	{
		interactions.GET("/", handlers.FindPeople)
		interactions.GET("/liked", handlers.ShowLiked)
		interactions.GET("/matches", handlers.ShowMatches)
		interactions.POST("", handlers.Interact)
		interactions.POST("/cancel", handlers.CancelInteraction)
	}

	hub := websocket.NewHub()
	go hub.StartServer(context.Background())
	
	v1.GET("/subscribe", middlewares.Auth(), handlers.Subscribe(hub))
}
