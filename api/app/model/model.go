package model

import (
	"time"

	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Base struct {
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:null" json:"deleted_at,omitempty"`
}

type User struct {
	Base
	ID      uint   `gorm:"unique;primaryKey;autoIncrement" json:"id"`
	Name    string `gorm:"unique" json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
	Status  bool   `json:"status"`
}

func (e *User) Disable() {
	e.Status = false
}

func (p *User) Enable() {
	p.Status = true
}

type Address struct {
	Base
	ID      uint   `gorm:"unique;primaryKey;autoIncrement" json:"id"`
	Street  string `json:"street"`
	City    string `json:"city"`
	Country string `json:"country"`
	User    *User  `gorm:"foreignKey:ID" json:"user,omitempty"`
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&User{}, &Address{})
	return db
}
