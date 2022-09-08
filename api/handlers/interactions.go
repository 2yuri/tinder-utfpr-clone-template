package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"tinderutf/db/repositories"
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
