package main

import (
	"encoding/json"

	"github.com/buurzx/cryptoexchange/orderbook"
	"github.com/labstack/echo/v4"
)

func main() {
	server := echo.New()

	exchange := NewExchange()

	server.POST("/orders", exchange.handlePlaceOrder)
	server.GET("/test", exchange.handleHeartBeat)

	server.Start(":3000")
}

type Market string

const (
	MarketETH Market = "ETH"
)

type OrderType string

const (
	LimitOrder  OrderType = "LIMIT"
	MarketOrder OrderType = "MARKET"
)

type Direction string

const (
	Ask Direction = "ASK"
	Bid Direction = "BID"
)

type Exchange struct {
	orderbooks map[Market]*orderbook.Orderbook
}

func NewExchange() *Exchange {
	orderbooks := make(map[Market]*orderbook.Orderbook)
	orderbooks[MarketETH] = orderbook.NewOrderBook()

	return &Exchange{
		orderbooks: orderbooks,
	}
}

type PlaceOrderRequest struct {
	Type      OrderType
	Direction Direction
	Size      float64 `json:",string"`
	Price     float64 `json:",string"`
	Market    Market
}

func (e *Exchange) handlePlaceOrder(c echo.Context) error {
	var placeOrderData PlaceOrderRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return c.JSON(500, map[string]any{"error": err.Error()})
	}

	market := Market(placeOrderData.Market)
	ob := e.orderbooks[market]

	order := orderbook.NewOrder(orderbook.OrderKind(placeOrderData.Direction), placeOrderData.Size)

	if placeOrderData.Type == MarketOrder {
		ob.PlaceMarketOrder(order)
	} else {
		ob.PlaceLimitOrder(placeOrderData.Price, order)
	}

	return c.JSON(200, map[string]any{"succes": true})
}

func (e *Exchange) handleHeartBeat(c echo.Context) error {
	return c.JSON(200, map[string]any{"heartbeat": "success"})
}
