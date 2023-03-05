package web

import (
	"encoding/json"

	"github.com/buurzx/cryptoexchange/src/entities"
	"github.com/labstack/echo/v4"
)

type placeOrder struct {
	orderbookRepo OrderbooksRepoIface
}

func NewPlaceOrderHandler(oredebookRepo OrderbooksRepoIface) *placeOrder {
	return &placeOrder{orderbookRepo: oredebookRepo}
}

func (p *placeOrder) Handle(c echo.Context) error {
	var placeOrderData entities.PlaceOrderRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return c.JSON(500, map[string]any{"error": []string{err.Error()}})
	}

	market := entities.Market(placeOrderData.Market)
	ob := p.orderbookRepo.FindByMarket(string(market))

	order := entities.NewOrder(entities.OrderKind(placeOrderData.Kind), placeOrderData.Size)

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
