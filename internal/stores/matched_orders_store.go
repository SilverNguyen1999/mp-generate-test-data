package stores

import (
	"mp-generate-test-data/internal/models"

	"gorm.io/gorm"
)

type MatchedOrdersStore struct {
	*gorm.DB
}

// for a service
func NewMatchedOrdersStore(db *gorm.DB) *MatchedOrdersStore {
	return &MatchedOrdersStore{db}
}

func (m *MatchedOrdersStore) Save(object *models.MatchedOrdersMd) error {
	object.BeforeCreate()
	return m.save(object)
}

func (m *MatchedOrdersStore) save(object *models.MatchedOrdersMd) error {
	return m.Create(object).Error
}

func (m *MatchedOrdersStore) InsertMany(slice []*models.MatchedOrdersMd) error {
	return m.Create(slice).Error
}
