package web

import (
	"encoding/json"
	"net/http"

	"github.com/buurzx/cryptoexchange/src/entities"
	"github.com/labstack/echo/v4"
)

type PlaceOrderRequest struct {
	Type   entities.OrderType
	Kind   entities.OrderKind
	Size   float64 `json:",string"`
	Price  float64 `json:",string"`
	Market entities.Market
}

type placeOrder struct {
	orderbookRepo OrderbooksRepoIface
}

func NewPlaceOrderHandler(oredebookRepo OrderbooksRepoIface) *placeOrder {
	return &placeOrder{orderbookRepo: oredebookRepo}
}

func (p *placeOrder) Handle(c echo.Context) error {
	var placeOrderData PlaceOrderRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return c.JSON(500, map[string]any{"error": []string{err.Error()}})
	}

	market := entities.Market(placeOrderData.Market)

	ob, err := p.orderbookRepo.FindByMarket(string(market))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]any{"errors": []string{"not_found"}})
	}

	order := entities.NewOrder(entities.OrderKind(placeOrderData.Kind), placeOrderData.Size, ob.ID)

	if placeOrderData.Type == entities.MarketOrder {
		matches := ob.PlaceMarketOrder(order)
		kind := order.Kind
		matchedOrders := []entities.MatchedOrder{}

		for _, match := range matches {
			var newMatch entities.MatchedOrder

			if kind == entities.Ask {
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
	if placeOrderData.Type == entities.LimitOrder {
		id := ob.PlaceLimitOrder(placeOrderData.Price, order)
		return c.JSON(200, map[string]any{"id": id})
	}

	return nil
}
