package handlers

import (
	"context"
	"net/http"
	"tinderutf/db/repositories"

	"log"
	"tinderutf/internal/validate"

	"github.com/gin-gonic/gin"
)

func FindPeople(c *gin.Context) {
	ctx := context.Background()

	id := c.GetString("userId")
	if id == "" {
		c.AbortWithStatusJSON(500, createDefaultError("invalid userId"))
		return
	}

	user, err := repositories.GetUserById(ctx, id)
	if err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	users, err := repositories.FindPeople(ctx, id, user.Options().Distance(), user.Info().Geo(), user.Options().SexPreference(), user.Sex())
	if err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	var response []interface{}
	for _, u := range users {
		response = append(response, u.ToJSON())
	}

	if len(response) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(200, createDefaultResponse(response))
}

func ShowLiked(c *gin.Context) {
	ctx := context.Background()

	id := c.GetString("userId")
	if id == "" {
		c.AbortWithStatusJSON(500, createDefaultError("invalid userId"))
		return
	}

	users, err := repositories.ShowLiked(ctx, id)
	if err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	var response []interface{}
	for _, u := range users {
		response = append(response, u.ToJSON())
	}

	if len(response) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(200, createDefaultResponse(response))
}

func ShowMatches(c *gin.Context) {
	ctx := context.Background()

	id := c.GetString("userId")
	if id == "" {
		c.AbortWithStatusJSON(500, createDefaultError("invalid userId"))
		return
	}

	users, err := repositories.ShowMatches(ctx, id)
	if err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	var response []interface{}
	for _, u := range users {
		response = append(response, u.ToJSON())
	}

	if len(response) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(200, createDefaultResponse(response))
}

func CreateInteraction(c *gin.Context) {
	ctx := context.Background()

	id := c.GetString("userId")
	if id == "" {
		c.AbortWithStatusJSON(400, createDefaultError("invalid userId"))
	}

	var req struct {
		Target string `validate:"required" json:"target"`
		Like bool `json:"like"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, createDefaultError(err.Error()))
		return
	}

	errors := validate.Struct(&req)
	if len(errors) > 0 {
		c.AbortWithStatusJSON(400, createDefaultError(errors...))
		return
	}

	if req.Target == id {
		c.AbortWithStatusJSON(400, createDefaultError("invalid targetId"))
		return
	}

	isLiked, err := repositories.CreateInteraction(ctx, id, req.Target, req.Like)
	if err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	if isLiked && req.Like {
		log.Println("é um match")
	}

	c.Status(200)
}

func CancelInteraction(c *gin.Context) {
	ctx := context.Background()

	id := c.GetString("userId")
	if id == "" {
		c.AbortWithStatusJSON(400, createDefaultError("invalid userId"))
	}

	var req struct {
		Target string `validate:"required" json:"target"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, createDefaultError(err.Error()))
		return
	}

	errors := validate.Struct(&req)
	if len(errors) > 0 {
		c.AbortWithStatusJSON(400, createDefaultError(errors...))
		return
	}

	if req.Target == id {
		c.AbortWithStatusJSON(400, createDefaultError("invalid targetId"))
		return
	}

	err = repositories.CancelInteraction(ctx, id, req.Target)
	if err != nil {
		c.AbortWithStatusJSON(400, createDefaultError(err.Error()))
		return
	}

	c.Status(200)
}