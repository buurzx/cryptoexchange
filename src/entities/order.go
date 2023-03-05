package entities

import (
	"math/rand"
	"time"
)

type OrderType string

const (
	LimitOrder  OrderType = "LIMIT"
	MarketOrder OrderType = "MARKET"
)

// OrderKind determine whether bid or ask type of order
type OrderKind string

const (
	Bid OrderKind = "BID"
	Ask OrderKind = "ASK"
)

type Orders []*Order

type Order struct {
	ID         int64
	Size       float64
	Kind       OrderKind
	Limit      *Limit
	Price      float64
	Timestampt int64
}

func NewOrder(orderKind OrderKind, size float64) *Order {
	return &Order{
		ID:         int64(rand.Intn(1000000000)),
		Size:       size,
		Kind:       orderKind,
		Timestampt: time.Now().UnixNano(),
	}
}

func (o *Order) IsFilled() bool {
	return o.Size == 0.0
}

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestampt > o[j].Timestampt }

type MatchedOrder struct {
	Price float64
	Size  float64
	ID    int64
}
