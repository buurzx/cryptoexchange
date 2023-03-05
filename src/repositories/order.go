// temporary realisation without db
package repositories

import (
	"github.com/buurzx/cryptoexchange/src/entities"
)

type OrderRepo struct {
	Records       map[entities.Market]*entities.Orderbook
	orderBookRepo *entities.Orderbook
}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{}
}

func (o *OrderRepo) FindByID(id int64) (order *entities.Order, err error) {
	order, exists := o.orderBookRepo.Orders[id]
	if !exists {
		return nil, entities.ErrNotFound
	}

	return
}
