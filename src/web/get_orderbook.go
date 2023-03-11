package web

import (
	"errors"
	"net/http"

	"github.com/buurzx/cryptoexchange/src/entities"
	"github.com/labstack/echo/v4"
)

type OrderbookData struct {
	TotalBidVolume float64
	TotalAskVolume float64
	Asks           []*entities.Order
	Bids           []*entities.Order
	OrderbookID    int64
}

type getOrderBook struct {
	orderbookRepo OrderbooksRepoIface
}

func NewGetOrderBookHandler(oredebookRepo OrderbooksRepoIface) *getOrderBook {
	return &getOrderBook{orderbookRepo: oredebookRepo}
}

func (p *getOrderBook) Handle(c echo.Context) error {
	market := entities.Market(c.Param("market"))

	ob, err := p.orderbookRepo.FindByMarket(string(market))
	if err != nil && errors.Is(entities.ErrNotFound, err) {
		return c.JSON(http.StatusNotFound, map[string]any{"errors": notFoundErrorResponse})
	}

	orderbookData := OrderbookData{
		TotalBidVolume: ob.BidTotalVolume(),
		TotalAskVolume: ob.AskTotalVolume(),
		Asks:           []*entities.Order{},
		Bids:           []*entities.Order{},
		OrderbookID:    ob.ID,
	}

	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			o := entities.Order{
				ID:          order.ID,
				Price:       limit.Price,
				Size:        order.Size,
				Kind:        order.Kind,
				Timestampt:  order.Timestampt,
				OrderbookID: order.OrderbookID,
			}

			orderbookData.Asks = append(orderbookData.Asks, &o)
		}
	}

	for _, limit := range ob.Bids() {
		for _, order := range limit.Orders {
			o := entities.Order{
				ID:          order.ID,
				Price:       limit.Price,
				Size:        order.Size,
				Kind:        order.Kind,
				Timestampt:  order.Timestampt,
				OrderbookID: order.OrderbookID,
			}

			orderbookData.Bids = append(orderbookData.Bids, &o)
		}
	}

	return c.JSON(http.StatusOK, orderbookData)
}
