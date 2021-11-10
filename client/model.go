package client

type Order struct {
	ID    string      `json:"id"`
	Items []OrderItem `json:"items"`
}

type OrderItem struct {
	Log Log `json:"log"`
}

type Log struct {
	Body string `json:"body"`
}
