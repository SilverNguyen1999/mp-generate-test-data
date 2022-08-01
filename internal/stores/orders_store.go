package stores

import (
	"mp-generate-test-data/internal/models"

	"gorm.io/gorm"
)

type OrdersStore struct {
	*gorm.DB
}

func NewOrdersStore(db *gorm.DB) *OrdersStore {
	return &OrdersStore{db}
}

func (m *OrdersStore) Save(object *models.OrdersMd) error {
	object.BeforeCreate()
	return m.save(object)
}

func (m *OrdersStore) save(object *models.OrdersMd) error {
	return m.Create(object).Error
}

func (m *OrdersStore) GetLargestOrderId() (*models.OrdersMd, bool, error) {
	var object = &models.OrdersMd{}

	err := m.Model(models.OrdersMd{}).Order("id DESC").Limit(1).First(object).Error
	if err == gorm.ErrRecordNotFound {
		return object, false, nil
	}
	return object, true, err
}

func (m *OrdersStore) InsertMany(slice []*models.OrdersMd) error {
	return m.Create(slice).Error
}
