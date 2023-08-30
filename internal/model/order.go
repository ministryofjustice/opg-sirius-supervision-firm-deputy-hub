package model

type OrderStatus struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type Orders struct {
	Order Order `json:"order"`
}

type Order struct {
	Id          int         `json:"id"`
	Client      Client      `json:"client"`
	OrderStatus OrderStatus `json:"orderStatus"`
}
