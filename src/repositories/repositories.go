// temporary realisation without db

package repositories

import (
	"github.com/buurzx/cryptoexchange/src/app"
	"github.com/buurzx/cryptoexchange/src/entities"
)

type Repo struct {
	Order     *OrderRepo
	Orderbook *OrderbookRepo
}

func NewRepos() *Repo {
	return &Repo{
		Orderbook: NewOrderbooksRepo(),
		Order:     NewOrderRepo(),
	}
}

func (r *Repo) Setup(a *app.Application) error {
	ob := entities.NewOrderBook()
	r.Orderbook.Records[entities.MarketETH] = ob

	r.Order.orderBookRepo = ob

	return nil
}

func (r *Repo) Run() error {
	return nil
}

func (r *Repo) Stop() error {
	return nil
}

const pluginName = "repositories"

func (r *Repo) Name() string {
	return pluginName
}
