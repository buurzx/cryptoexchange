package orderbook

import (
	"fmt"
	"sort"
	"time"
)

// Match finded match to ask or bid
type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

// OrderKind determine whether bid or ask type of order
type OrderKind int

const (
	Bid OrderKind = 0
	Ask OrderKind = 1
)

// Order
type Order struct {
	Size       float64
	Kind       OrderKind
	Limit      *Limit
	Timestampt int64
}

type Orders []*Order

func NewOrder(orderKind OrderKind, size float64) *Order {
	return &Order{
		Size:       size,
		Kind:       orderKind,
		Timestampt: time.Now().UnixNano(),
	}
}

func (o *Order) IsFilled() bool {
	return o.Size == 0.0
}

// func (o *Order) String() string {
// 	return fmt.Sprintf("[size: %.2f, ]", o.Size)
// }

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestampt > o[j].Timestampt }

// Orderbook holds a bunch of asks and bids
type Orderbook struct {
	asks []*Limit
	bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

// Limit holds a bunch of orders on certain price level
type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit

type ByBestAsk struct {
	Limits
}

func (a ByBestAsk) Len() int           { return len(a.Limits) }
func (a ByBestAsk) Swap(i, j int)      { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a ByBestAsk) Less(i, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }

type ByBestBid struct{ Limits }

func (b ByBestBid) Len() int           { return len(b.Limits) }
func (b ByBestBid) Swap(i, j int)      { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
func (b ByBestBid) Less(i, j int) bool { return b.Limits[i].Price > b.Limits[j].Price }

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

func (l *Limit) String(o *Order) string {
	return fmt.Sprintf("[price: %.2f | totalVolue: %.2f", l.Price, l.TotalVolume)
}

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (l *Limit) DeleteOrder(o *Order) {
	// need to fill gap from the deleting order
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			l.Orders[i] = l.Orders[len(l.Orders)-1]
			l.Orders = l.Orders[:len(l.Orders)-1]
		}
	}
	// destroy links
	o.Limit = nil
	l.TotalVolume -= o.Size

	sort.Sort(l.Orders)
}

func (l *Limit) Fill(o *Order) []Match {
	var (
		ordersToDelete []*Order
		matches        []Match
	)

	for _, order := range l.Orders {
		match := l.fillOrder(order, o)
		matches = append(matches, match)

		l.TotalVolume -= match.SizeFilled

		if order.IsFilled() {
			ordersToDelete = append(ordersToDelete, order)
			break
		}
	}

	for _, order := range ordersToDelete {
		l.DeleteOrder(order)
	}

	return matches
}

func (l *Limit) fillOrder(order1, order2 *Order) Match {
	var (
		bid        *Order
		ask        *Order
		sizeFilled float64
	)

	if order1.Kind == Bid {
		bid = order1
		ask = order2
	} else {
		bid = order2
		ask = order1
	}

	if order1.Size > order2.Size {
		order1.Size -= order2.Size
		sizeFilled = order2.Size
		order2.Size = 0.0
	} else {
		order2.Size -= order1.Size
		sizeFilled = order1.Size
		order1.Size = 0.0
	}

	return Match{
		Bid:        bid,
		Ask:        ask,
		SizeFilled: sizeFilled,
		Price:      l.Price,
	}
}

func NewOrderBook() *Orderbook {
	return &Orderbook{
		asks:      []*Limit{},
		bids:      []*Limit{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

// PlaceMarketOrder places market order which drained limit orders from current price and lower
func (ob *Orderbook) PlaceMarketOrder(o *Order) []Match {
	matches := []Match{}

	if o.Kind == Bid {
		if o.Size > ob.AskTotalVolume() {
			panic(fmt.Errorf("not enough volume [size: %.2f] for market order [size: %.2f]", ob.AskTotalVolume(), o.Size))
		}
		for _, limit := range ob.Asks() {
			limitMatches := limit.Fill(o)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				ob.clearLimit(Ask, limit)
			}
		}
	} else {
		if o.Size > ob.BidTotalVolume() {
			panic(fmt.Errorf("not enough volume [size: %.2f] for market order [size: %.2f]", ob.BidTotalVolume(), o.Size))
		}

		for _, limit := range ob.Bids() {
			limitMatches := limit.Fill(o)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				ob.clearLimit(Bid, limit)
			}
		}
	}

	return matches
}

// PlaceLimitOrder adds limit order to order book
func (ob *Orderbook) PlaceLimitOrder(price float64, o *Order) {
	var limit *Limit

	if o.Kind == Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AskLimits[price]
	}

	if limit == nil {
		limit = NewLimit(price)

		if o.Kind == Bid {
			ob.bids = append(ob.bids, limit)
		} else {
			ob.asks = append(ob.asks, limit)
		}
	}

	limit.AddOrder(o)
}

func (ob *Orderbook) clearLimit(kind OrderKind, l *Limit) {
	if kind == Bid {
		delete(ob.BidLimits, l.Price)

		for i := 0; i < len(ob.bids); i++ {
			if ob.bids[i] == l {
				ob.bids[i] = ob.bids[len(ob.bids)-1]
				ob.bids = ob.bids[:len(ob.bids)-1]
			}
		}
	} else {
		delete(ob.AskLimits, l.Price)

		for i := 0; i < len(ob.asks); i++ {
			if ob.asks[i] == l {
				ob.asks[i] = ob.asks[len(ob.asks)-1]
				ob.asks = ob.asks[:len(ob.asks)-1]
			}
		}
	}
}

func (ob *Orderbook) CancelOrder(o *Order) {
	limit := o.Limit
	limit.DeleteOrder(o)
}

func (ob *Orderbook) BidTotalVolume() float64 {
	totalVolume := 0.0

	for i := 0; i < len(ob.bids); i++ {
		totalVolume += ob.bids[i].TotalVolume
	}

	return totalVolume
}

func (ob *Orderbook) AskTotalVolume() float64 {
	totalVolume := 0.0

	for i := 0; i < len(ob.asks); i++ {
		totalVolume += ob.asks[i].TotalVolume
	}

	return totalVolume
}

func (ob *Orderbook) Asks() []*Limit {
	sort.Sort(ByBestAsk{ob.asks})
	return ob.asks
}

func (ob *Orderbook) Bids() []*Limit {
	sort.Sort(ByBestBid{ob.bids})
	return ob.bids
}
