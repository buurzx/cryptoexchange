package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/buurzx/cryptoexchange/orderbook"
	"github.com/labstack/echo/v4"
)

func main() {
	server := echo.New()

	exchange := NewExchange()

	server.POST("/orders", exchange.handlePlaceOrder)
	server.GET("/orderbook/:market", exchange.handleGetOrderbook)
	server.DELETE("/orders/:id", exchange.handleCancelOrder)

	server.GET("/heartbeat", exchange.handleHeartBeat)

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

type Kind string

const (
	Ask Kind = "ASK"
	Bid Kind = "BID"
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
	Type   OrderType
	Kind   Kind
	Size   float64 `json:",string"`
	Price  float64 `json:",string"`
	Market Market
}

type Order struct {
	ID        int64
	Size      float64
	Timestamp int64
	Price     float64
	Kind      string
}

type OrderbookData struct {
	TotalBidVolume float64
	TotalAskVolume float64
	Asks           []*Order
	Bids           []*Order
}

type MatchedOrder struct {
	Price float64
	Size  float64
	ID    int64
}

func (e *Exchange) handlePlaceOrder(c echo.Context) error {
	var placeOrderData PlaceOrderRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return c.JSON(500, map[string]any{"error": []string{err.Error()}})
	}

	market := Market(placeOrderData.Market)
	ob := e.orderbooks[market]

	order := orderbook.NewOrder(orderbook.OrderKind(placeOrderData.Kind), placeOrderData.Size)

	if placeOrderData.Type == MarketOrder {
		matches := ob.PlaceMarketOrder(order)
		kind := order.Kind
		matchedOrders := []MatchedOrder{}

		for _, match := range matches {
			var newMatch MatchedOrder

			if kind == orderbook.Ask {
				newMatch.ID = match.Ask.ID
			} else {
				newMatch.ID = match.Bid.ID
			}
			newMatch.Size = match.SizeFilled
			newMatch.Price = match.Price

			matchedOrders = append(matchedOrders, newMatch)
		}

		return c.JSON(200, map[string]any{"matches": matchedOrders})
	}
	if placeOrderData.Type == LimitOrder {
		id := ob.PlaceLimitOrder(placeOrderData.Price, order)
		return c.JSON(200, map[string]any{"id": id})
	}

	return nil
}

func (e *Exchange) handleHeartBeat(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"heartbeat": "success"})
}

func (e *Exchange) handleGetOrderbook(c echo.Context) error {
	market := Market(c.Param("market"))

	ob, ok := e.orderbooks[market]
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]any{"errors": []string{"orderbook not found"}})
	}

	orderbookData := OrderbookData{
		TotalBidVolume: ob.BidTotalVolume(),
		TotalAskVolume: ob.AskTotalVolume(),
		Asks:           []*Order{},
		Bids:           []*Order{},
	}

	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			o := Order{
				ID:        order.ID,
				Price:     limit.Price,
				Size:      order.Size,
				Kind:      string(order.Kind),
				Timestamp: order.Timestampt,
			}

			orderbookData.Asks = append(orderbookData.Asks, &o)
		}
	}

	for _, limit := range ob.Bids() {
		for _, order := range limit.Orders {
			o := Order{
				ID:        order.ID,
				Price:     limit.Price,
				Size:      order.Size,
				Kind:      string(order.Kind),
				Timestamp: order.Timestampt,
			}

			orderbookData.Bids = append(orderbookData.Bids, &o)
		}
	}

	return c.JSON(http.StatusOK, orderbookData)
}

func (e *Exchange) handleCancelOrder(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"errors": []string{err.Error()}})
	}

	ob := e.orderbooks[MarketETH]
	order, exists := ob.Orders[int64(id)]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]any{"errors": []string{"Order not found"}})
	}

	ob.CancelOrder(order)

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}
