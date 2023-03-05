package entities

type Market string

const (
	MarketETH Market = "ETH"
)

type Exchange struct {
	Orderbooks map[Market]*Orderbook
}

func NewExchange() *Exchange {
	orderbooks := make(map[Market]*Orderbook)
	orderbooks[MarketETH] = NewOrderBook()

	return &Exchange{
		Orderbooks: orderbooks,
	}
}
