package stores

import (
	"mp-generate-test-data/internal/models"

	"gorm.io/gorm"
)

type OrderAssetsStore struct {
	*gorm.DB
}

func NewOrderAssetsStore(db *gorm.DB) *OrderAssetsStore {
	return &OrderAssetsStore{db}
}

func (m *OrderAssetsStore) Save(object *models.OrderAssetsMd) error {
	object.BeforeCreate()
	return m.save(object)
}

func (m *OrderAssetsStore) save(object *models.OrderAssetsMd) error {
	return m.Create(object).Error
}

func (m *OrderAssetsStore) InsertMany(slice []*models.OrderAssetsMd) error {
	return m.Create(slice).Error
}
