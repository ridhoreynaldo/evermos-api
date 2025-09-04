package models

import (
    "time"
)

type DetailTrx struct {
    ID         uint    `json:"id" gorm:"primaryKey"`
    IDTrx      uint    `json:"id_trx" gorm:"not null"`
    IDLogProduk uint   `json:"id_log_produk" gorm:"not null"`
    IDToko     uint    `json:"id_toko" gorm:"not null"`
    Kuantitas  int     `json:"kuantitas" gorm:"not null"`
    HargaTotal int     `json:"harga_total" gorm:"not null"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    
    // Relations
    Trx        Trx        `json:"trx,omitempty" gorm:"foreignKey:IDTrx"`
    LogProduct LogProduct `json:"log_product,omitempty" gorm:"foreignKey:IDLogProduk"`
    Toko       Toko       `json:"toko,omitempty" gorm:"foreignKey:IDToko"`
}