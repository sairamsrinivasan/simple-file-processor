package db

import (
	"gorm.io/gorm"
)

// GormDB is a wrapper interface for gorm.DB
// So we can use dependency injection for testing purposes
type GormDB interface {
	Create(interface{}) *gorm.DB
	AutoMigrate(...interface{}) error
	Model(value interface{}) *gorm.DB
}

type gormDB struct {
	db *gorm.DB
}

// NewGormDB creates a new gormDB instance with the given gorm.DB
// the gorm instance is passed as an interface to allow for mocking in tests
func NewGormDB(db *gorm.DB) GormDB {
	return &gormDB{db: db}
}

func (gdb gormDB) Create(value interface{}) *gorm.DB {
	return gdb.db.Create(value)
}

func (gdb gormDB) AutoMigrate(value ...interface{}) error {
	return gdb.db.AutoMigrate(value...)
}

func (gdb gormDB) Model(value interface{}) *gorm.DB {
	return gdb.db.Model(value)
}
