package web

import (
	"net/http"

	"github.com/buurzx/cryptoexchange/src/entities"
	"github.com/labstack/echo/v4"
)

type OrderbookData struct {
	TotalBidVolume float64
	TotalAskVolume float64
	Asks           []*entities.Order
	Bids           []*entities.Order
}

type getOrder struct {
	orderbookRepo OrderbooksRepoIface
}

func NewGetOrderHandler(oredebookRepo OrderbooksRepoIface) *getOrder {
	return &getOrder{orderbookRepo: oredebookRepo}
}

func (p *getOrder) Handle(c echo.Context) error {
	market := entities.Market(c.Param("market"))

	ob := p.orderbookRepo.FindByMarket(string(market))
	if ob == nil {
		return c.JSON(http.StatusNotFound, map[string]any{"errors": []string{"orderbook not found"}})
	}

	orderbookData := OrderbookData{
		TotalBidVolume: ob.BidTotalVolume(),
		TotalAskVolume: ob.AskTotalVolume(),
		Asks:           []*entities.Order{},
		Bids:           []*entities.Order{},
	}

	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			o := entities.Order{
				ID:         order.ID,
				Price:      limit.Price,
				Size:       order.Size,
				Kind:       order.Kind,
				Timestampt: order.Timestampt,
			}

			orderbookData.Asks = append(orderbookData.Asks, &o)
		}
	}

	for _, limit := range ob.Bids() {
		for _, order := range limit.Orders {
			o := entities.Order{
				ID:         order.ID,
				Price:      limit.Price,
				Size:       order.Size,
				Kind:       order.Kind,
				Timestampt: order.Timestampt,
			}

			orderbookData.Bids = append(orderbookData.Bids, &o)
		}
	}

	return c.JSON(http.StatusOK, orderbookData)
}
