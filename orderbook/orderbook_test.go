package orderbook

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrder1 := NewOrder(Bid, 5)
	buyOrder2 := NewOrder(Bid, 8)
	buyOrder3 := NewOrder(Bid, 12)

	l.AddOrder(buyOrder1)
	l.AddOrder(buyOrder2)
	l.AddOrder(buyOrder3)

	l.DeleteOrder(buyOrder2)

	require.Equal(t, 17.0, l.TotalVolume)

	// fmt.Println(l.TotalVolume)
	// fmt.Printf("Limit %#v \n", l)
}

func TestPlaceLimitOrder(t *testing.T) {
	price := 10_000.0
	ob := NewOrderBook()

	sellOrder1 := NewOrder(Ask, 10)
	sellOrder2 := NewOrder(Ask, 8)
	ob.PlaceLimitOrder(price, sellOrder1)
	ob.PlaceLimitOrder(price-1000.0, sellOrder2)

	assert(t, len(ob.asks), 2)

	buyOrder1 := NewOrder(Bid, 10)
	buyOrder2 := NewOrder(Bid, 12)
	ob.PlaceLimitOrder(price, buyOrder1)
	ob.PlaceLimitOrder(price, buyOrder2)

	assert(t, len(ob.bids), 2)

	assert(t, ob.Orders[sellOrder1.ID], sellOrder1)
	assert(t, ob.Orders[buyOrder1.ID], buyOrder1)
}

func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrderBook()

	sellOrder := NewOrder(Ask, 20)
	ob.PlaceLimitOrder(10_000, sellOrder)

	buyOrder := NewOrder(Bid, 10)
	matches := ob.PlaceMarketOrder(buyOrder)

	assert(t, len(matches), 1)
	assert(t, len(ob.asks), 1)
	assert(t, ob.AskTotalVolume(), 10.0)

	assert(t, matches[0].Ask, sellOrder)
	assert(t, matches[0].Bid, buyOrder)
	assert(t, matches[0].SizeFilled, 10.0)
	assert(t, matches[0].Price, 10_000.0)
	assert(t, buyOrder.IsFilled(), true)

	fmt.Printf("%+v", matches)
}

func TestPlaceMarketOrderMultiFillAsk(t *testing.T) {
	ob := NewOrderBook()

	buyOrder1 := NewOrder(Bid, 5)
	buyOrder2 := NewOrder(Bid, 8)
	buyOrder3 := NewOrder(Bid, 10)
	buyOrder4 := NewOrder(Bid, 1)

	ob.PlaceLimitOrder(5_000, buyOrder3)
	ob.PlaceLimitOrder(5_000, buyOrder4)
	ob.PlaceLimitOrder(9_000, buyOrder2)
	ob.PlaceLimitOrder(10_000, buyOrder1)

	assert(t, ob.BidTotalVolume(), 24.00)

	sellOrder := NewOrder(Ask, 20)
	matches := ob.PlaceMarketOrder(sellOrder)

	assert(t, ob.BidTotalVolume(), 4.0)
	assert(t, len(matches), 4)
	assert(t, len(ob.bids), 2)

	fmt.Printf("%+v", matches)

}

func TestPlaceMarketOrderMultiFillBid(t *testing.T) {
	ob := NewOrderBook()

	sellOrder1 := NewOrder(Ask, 5)
	sellOrder2 := NewOrder(Ask, 8)
	sellOrder3 := NewOrder(Ask, 10)

	ob.PlaceLimitOrder(5_000, sellOrder3)
	ob.PlaceLimitOrder(9_000, sellOrder2)
	ob.PlaceLimitOrder(10_000, sellOrder1)

	assert(t, ob.AskTotalVolume(), 23.00)

	buyOrder := NewOrder(Bid, 20)
	matches := ob.PlaceMarketOrder(buyOrder)

	assert(t, ob.AskTotalVolume(), 3.0)
	assert(t, len(ob.asks), 1)
	assert(t, len(matches), 3)

	fmt.Printf("%+v", matches)

}

func TestCancelOrder(t *testing.T) {
	ob := NewOrderBook()

	buyOrder := NewOrder(Bid, 4)
	ob.PlaceLimitOrder(10_000.0, buyOrder)

	assert(t, ob.BidTotalVolume(), 4.0)

	ob.CancelOrder(buyOrder)

	assert(t, ob.BidTotalVolume(), 0.0)

	_, exists := ob.Orders[buyOrder.ID]
	assert(t, exists, false)
}
