package stores

import "gorm.io/gorm"

type MainStore struct {
	*gorm.DB
}

// NewMainStore creates a new instance of MainStore, contains whole common functions
// for a service

func NewMainStore(db *gorm.DB) *MainStore {
	return &MainStore{db}
}

func (m *MainStore) RelationalDatabaseCheck() error {
	return m.Raw("SELECT 1").Error
}
