package repositories

import (
	"github.com/buurzx/cryptoexchange/src/app"
	"github.com/buurzx/cryptoexchange/src/entities"
)

type OrderbookRepo struct {
	Records map[entities.Market]*entities.Orderbook
}

const pluginName = "orderbookRepo"

func NewOrderbooksRepo() *OrderbookRepo {
	return &OrderbookRepo{
		Records: map[entities.Market]*entities.Orderbook{},
	}
}

func (ob *OrderbookRepo) FindByMarket(market string) *entities.Orderbook {
	return ob.Records[entities.Market(market)]
}

func (ob *OrderbookRepo) Setup(a *app.Application) error {
	ob.Records[entities.MarketETH] = entities.NewOrderBook()

	return nil
}

func (ob *OrderbookRepo) Run() error {
	return nil
}

func (ob *OrderbookRepo) Stop() error {
	return nil
}

func (ob *OrderbookRepo) Name() string {
	return pluginName
}
