package models

type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume float64
}
