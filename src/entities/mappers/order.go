package mappers

import (
	"github.com/buurzx/cryptoexchange/src/entities"
	"github.com/buurzx/cryptoexchange/src/repositories/models"
)

func OrderModelToEntity(o models.Order) *entities.Order {
	return &entities.Order{
		ID:         o.ID,
		Size:       o.Size,
		Kind:       entities.OrderKind(o.Kind),
		Timestampt: o.Timestampt,
		Price:      o.Price,
	}
}
