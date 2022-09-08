package models

import (
	"github.com/shopspring/decimal"
	"time"
	"tinderutf/domain"
)

type User struct {
	Id            int32
	Name          string
	BirthDate     time.Time
	Email         string
	Instagram     string
	Password      string
	About         string
	IsActive      bool
	Sex           domain.Sex
	SexPreference domain.Sex
	FindDistance  int
	Latitude      decimal.NullDecimal
	Longitude     decimal.NullDecimal
	ExternalId    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
