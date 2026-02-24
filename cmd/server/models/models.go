package models

type Item struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Picked bool   `json:"picked"`
	SeqNum int    `json:"seqnum"`
	Price  int    `json:"price"`
}

type ItemUpdate struct {
	Picked bool `json:"picked"`
}

type Email struct {
	EmailAddress string `json:"emailAddress"`
}

type ListDetails struct {
	TotalPrice    int `json:"totalprice"`
	SpendingLimit int `json:"spendingLimit"`
}
