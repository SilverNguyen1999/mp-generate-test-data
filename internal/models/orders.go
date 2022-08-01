package models

type Assets struct {
	Erc      int8   `json:"erc"`     // 0 , 1 , 2 --> 0: erc20 ...
	Address  string `json:"address"` // 4 contract address (verified contract)
	Id       int64  `json:"id"`      // exam: id of axie
	Quantity int    `json:"quantity"`
}

type OrdersMd struct {
	Id                  int64   `json:"id" gorm:"primaryKey"`
	Maker               string  `json:"maker"`
	Kind                int8    `json:"kind"`
	Assets              string  `json:"assets"`
	ExpiredAt           int64   `json:"expired_at"`
	TokenPayment        string  `json:"token_payment"`
	StartedAt           int64   `json:"started_at"`
	BasePrice           float64 `json:"base_price"`
	EndedAt             int64   `json:"ended_at"`
	EndedPrice          float64 `json:"ended_price"`
	ExpectedState       string  `json:"expected_state"`
	Nonce               int64   `json:"nonce"`
	MarketFeePercentage uint32  `json:"market_fee_percentage"`
	Signature           string  `json:"signature"`
	Hash                string  `json:"hash"`
}

func (m *OrdersMd) BeforeCreate() {

}

func (m *OrdersMd) TableName() string {
	return "orders"
}
