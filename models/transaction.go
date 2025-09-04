package models

import (
    "time"
)

type Trx struct {
    ID            uint        `json:"id" gorm:"primaryKey"`
    IDUser        uint        `json:"id_user" gorm:"not null"`
    AlamatKirim   uint        `json:"alamat_kirim" gorm:"not null"`
    HargaTotal    int         `json:"harga_total" gorm:"not null"`
    KodeInvoice   string      `json:"kode_invoice" gorm:"unique;not null"`
    MetodeBayar   string      `json:"metode_bayar" gorm:"not null"`
    CreatedAt     time.Time   `json:"created_at"`
    UpdatedAt     time.Time   `json:"updated_at"`
    
    // Relations
    User          User        `json:"user,omitempty" gorm:"foreignKey:IDUser"`
    Alamat        Alamat      `json:"alamat,omitempty" gorm:"foreignKey:AlamatKirim"`
    DetailTrx     []DetailTrx `json:"detail_trx,omitempty" gorm:"foreignKey:IDTrx"`
}