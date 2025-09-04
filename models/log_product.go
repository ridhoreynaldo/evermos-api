package models

import (
    "time"
)

type LogProduct struct {
    ID             uint      `json:"id" gorm:"primaryKey"`
    IDProduk       uint      `json:"id_produk" gorm:"not null"`
    NamaProduk     string    `json:"nama_produk" gorm:"not null"`
    Slug           string    `json:"slug" gorm:"not null"`
    HargaReseller  int       `json:"harga_reseller" gorm:"not null"`
    HargaKonsumen  int       `json:"harga_konsumen" gorm:"not null"`
    Deskripsi      string    `json:"deskripsi"`
    IDToko         uint      `json:"id_toko" gorm:"not null"`
    CategoryID     uint      `json:"category_id" gorm:"not null"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
    
    // Relations
    Product        Product   `json:"product,omitempty" gorm:"foreignKey:IDProduk"`
}