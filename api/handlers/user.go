package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"path/filepath"
	"time"
	"tinderutf/db/repositories"
	"tinderutf/domain"
	"tinderutf/internal/auth"
	"tinderutf/internal/validate"
)

func createDefaultResponse(data interface{}) interface{} {
	return &gin.H{
		"data": data,
	}
}

func createDefaultError(data ...string) interface{} {
	return &gin.H{
		"error": data,
	}
}

func Login(c *gin.Context) {
	ctx := context.Background()

	var req struct {
		Email    string `validate:"required,email" json:"email"`
		Password string `validate:"required,min=6" json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	if errors := validate.Struct(&req); len(errors) > 0 {
		c.AbortWithStatusJSON(500, createDefaultError(errors...))
		return
	}

	user, err := repositories.GetUserByEmail(ctx, req.Email)
	if err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	if err := user.ComparePassword(req.Password); err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	token, err := auth.JWT.GenerateToken(user.Id())
	if err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	c.JSON(200, createDefaultResponse(gin.H{
		"token": token,
	}))
	return
}

func Me(c *gin.Context) {
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

	c.JSON(200, createDefaultResponse(user.ToJSON()))
	return
}

func CreateUser(c *gin.Context) {
	ctx := context.Background()

	var req struct {
		Name      string     `validate:"required" json:"name"`
		Email     string     `validate:"required,email" json:"email"`
		BirthDate time.Time  `validate:"required" json:"birth_date"`
		Password  string     `validate:"required,min=6" json:"password"`
		Sex       domain.Sex `validate:"required" json:"sex"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	if errors := validate.Struct(&req); len(errors) > 0 {
		c.AbortWithStatusJSON(500, createDefaultError(errors...))
		return
	}

	user, err := domain.NewUser(req.Name, req.Email, req.Password, req.BirthDate, req.Sex)
	if err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	id, err := repositories.CreateUser(ctx, user)
	if err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	c.JSON(200, createDefaultResponse(gin.H{
		"id": id,
	}))
	return
}

func SetGeolocation(c *gin.Context) {
	ctx := context.Background()

	id := c.GetString("userId")
	if id == "" {
		c.AbortWithStatusJSON(500, createDefaultError("invalid userId"))
		return
	}

	var req struct {
		Latitude  decimal.Decimal `validate:"required" json:"latitude"`
		Longitude decimal.Decimal `validate:"required" json:"longitude"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	if errors := validate.Struct(&req); len(errors) > 0 {
		c.AbortWithStatusJSON(500, createDefaultError(errors...))
		return
	}

	if err := repositories.SetUserLocation(ctx, req.Latitude, req.Longitude, id); err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	c.Status(200)
	return
}

func SetCustomization(c *gin.Context) {
	ctx := context.Background()

	id := c.GetString("userId")
	if id == "" {
		c.AbortWithStatusJSON(500, createDefaultError("invalid userId"))
		return
	}

	var req struct {
		Instagram string     `validate:"required" json:"instagram"`
		About     string     `validate:"required" json:"about"`
		Distance  int        `validate:"required" json:"distance"`
		SexPref   domain.Sex `validate:"required" json:"sex_preference"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	if errors := validate.Struct(&req); len(errors) > 0 {
		c.AbortWithStatusJSON(500, createDefaultError(errors...))
		return
	}

	if err := repositories.SetUserCustomization(ctx, req.Instagram, req.About, req.Distance, req.SexPref, id); err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	c.Status(200)
	return
}

func SetProfileImage(c *gin.Context) {
	id := c.GetString("userId")
	if id == "" {
		c.AbortWithStatusJSON(500, createDefaultError("invalid userId"))
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.AbortWithStatusJSON(400, createDefaultError("missing image file"))
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpeg" {
		c.AbortWithStatusJSON(http.StatusBadRequest, createDefaultError("file need to be .jpeg"))
		return
	}

	if err := c.SaveUploadedFile(file, "./images/"+id+ext); err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}
}

func GetById(c *gin.Context) {
	ctx := context.Background()

	user, err := repositories.GetUserById(ctx, c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(500, createDefaultError(err.Error()))
		return
	}

	c.JSON(200, createDefaultResponse(user.ToJSON()))
	return
}
