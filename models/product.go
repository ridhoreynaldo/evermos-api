package models

import (
    "time"
)

type Product struct {
    ID             uint      `json:"id" gorm:"primaryKey"`
    NamaProduk     string    `json:"nama_produk" gorm:"not null" validate:"required"`
    Slug           string    `json:"slug" gorm:"unique;not null"`
    HargaReseller  int       `json:"harga_reseller" gorm:"not null" validate:"required,min=0"`
    HargaKonsumen  int       `json:"harga_konsumen" gorm:"not null" validate:"required,min=0"`
    Stok           int       `json:"stok" gorm:"not null" validate:"required,min=0"`
    Deskripsi      string    `json:"deskripsi"`
    URLFoto        []string  `json:"url_foto" gorm:"type:json"`
    IDToko         uint      `json:"id_toko" gorm:"not null"`
    CategoryID     uint      `json:"category_id" gorm:"not null"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
    
    // Relations
    Toko           Toko      `json:"toko,omitempty" gorm:"foreignKey:IDToko"`
    Category       Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}