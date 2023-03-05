package models

// Orderbook holds a bunch of asks and bids
type Orderbook struct {
	asks []*Limit
	bids []*Limit

	Orders map[int64]*Order

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}
