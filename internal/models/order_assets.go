package models

type OrderAssetsMd struct {
	OrderId  int64  `json:"order_id"`
	Erc      int8   `json:"erc"`
	Address  string `json:"address"`
	Id       int64  `json:"id"`
	Quantity int64  `json:"quantity"`
}

func (m *OrderAssetsMd) BeforeCreate() {

}

func (m *OrderAssetsMd) TableName() string {
	return "order_assets"
}
