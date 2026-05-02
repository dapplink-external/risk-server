package database

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Addresses struct {
	Guid        string    `gorm:"primaryKey; column:guid"`
	Address     string    `gorm:"column:address"`
	AddressType string    `gorm:"column:address_type"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

type AddressesView interface {
	QueryAddresses(string) ([]*Addresses, error)
}

type AddressesDB interface {
	AddressesView

	StoreAddresses(string, []*Addresses) error
}

type addressesDB struct {
	gorm *gorm.DB
}

func NewAddressesDB(db *gorm.DB) AddressesDB {
	return &addressesDB{gorm: db}
}

func (db *addressesDB) StoreAddresses(businessId string, addresses []*Addresses) error {
	result := db.gorm.Table("addresses_"+businessId).CreateInBatches(&addresses, len(addresses))
	return result.Error
}

func (db *addressesDB) QueryAddresses(requestId string) ([]*Addresses, error) {
	var addressListEntry []*Addresses
	err := db.gorm.Table("addresses_" + requestId).Find(&addressListEntry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return addressListEntry, nil
}
