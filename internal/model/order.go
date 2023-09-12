package model

type OrderStatus struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type Order struct {
	Id          int         `json:"id"`
	Client      Client      `json:"client"`
	OrderStatus OrderStatus `json:"orderStatus"`
}

type Orders struct {
	Order Order `json:"order"`
}
