package models

import (
    "time"
)

type Toko struct {
    ID       uint      `json:"id" gorm:"primaryKey"`
    NamaToko string    `json:"nama_toko" gorm:"not null" validate:"required"`
    UrlFoto  string    `json:"url_foto"`
    IDUser   uint      `json:"id_user" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    // Relations
    User     User      `json:"user,omitempty" gorm:"foreignKey:IDUser"`
    Products []Product `json:"products,omitempty" gorm:"foreignKey:IDToko"`
}