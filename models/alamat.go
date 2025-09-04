package models

import (
    "time"
)

type Alamat struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    JudulAlamat  string    `json:"judul_alamat" gorm:"not null" validate:"required"`
    NamaPenerima string    `json:"nama_penerima" gorm:"not null" validate:"required"`
    NoTelp       string    `json:"no_telp" gorm:"not null" validate:"required"`
    DetailAlamat string    `json:"detail_alamat" gorm:"not null" validate:"required"`
    IDUser       uint      `json:"id_user" gorm:"not null"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    
    // Relations
    User         User      `json:"user,omitempty" gorm:"foreignKey:IDUser"`
}