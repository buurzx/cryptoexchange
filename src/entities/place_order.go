package entities

type PlaceOrderRequest struct {
	Type   OrderType
	Kind   OrderKind
	Size   float64 `json:",string"`
	Price  float64 `json:",string"`
	Market Market
}
