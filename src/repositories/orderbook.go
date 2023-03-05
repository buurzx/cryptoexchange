// temporary realisation without db

package repositories

import (
	"github.com/buurzx/cryptoexchange/src/entities"
)

type OrderbookRepo struct {
	Records map[entities.Market]*entities.Orderbook
}

func NewOrderbooksRepo() *OrderbookRepo {
	return &OrderbookRepo{
		Records: map[entities.Market]*entities.Orderbook{},
	}
}

func (ob *OrderbookRepo) FindByMarket(market string) (*entities.Orderbook, error) {
	result, exists := ob.Records[entities.Market(market)]
	if !exists {
		return nil, entities.ErrNotFound
	}

	return result, nil
}
func (ob *OrderbookRepo) FindByID(id int64) (*entities.Orderbook, error) {
	for _, ob := range ob.Records {
		if ob.ID == id {
			return ob, nil
		}
	}

	return nil, entities.ErrNotFound
}
