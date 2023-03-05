package models

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
