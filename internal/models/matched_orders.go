package models

type MatchedOrdersMd struct {
	OrderId     int64  `json:"order_id"`
	Price       int64  `json:"price"`
	Matcher     string `json:"matcher"`
	BlockNumber int64  `json:"block_number"`
	TxHash      string `json:"tx_hash"`
}

func init() {

}

func (m *MatchedOrdersMd) BeforeCreate() {

}

func (m MatchedOrdersMd) TableName() string {
	return "matched_orders"
}
