package domain

import (
	"github.com/shopspring/decimal"
)

type Geo struct {
	latitude  decimal.Decimal
	longitude decimal.Decimal
}

func (g *Geo) Latitude() decimal.Decimal {
	return g.latitude
}

func (g *Geo) Longitude() decimal.Decimal {
	return g.longitude
}

func NewGeo(latitude, longitude decimal.Decimal) *Geo {
	return &Geo{latitude: latitude, longitude: longitude}
}
