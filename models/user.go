package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Nama        string         `json:"nama" gorm:"not null" validate:"required"`
    KataSandi   string         `json:"kata_sandi,omitempty" gorm:"not null" validate:"required,min=6"`
    NoTelp      string         `json:"no_telp" gorm:"unique;not null" validate:"required"`
    TanggalLahir time.Time     `json:"tanggal_lahir" validate:"required"`
    JenisKelamin string        `json:"jenis_kelamin" gorm:"type:enum('L','P')" validate:"required,oneof=L P"`
    Tentang     string         `json:"tentang"`
    Pekerjaan   string         `json:"pekerjaan"`
    Email       string         `json:"email" gorm:"unique;not null" validate:"required,email"`
    IDProvinsi  string         `json:"id_provinsi" validate:"required"`
    IDKota      string         `json:"id_kota" validate:"required"`
    IsAdmin     bool           `json:"is_admin" gorm:"default:false"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
    
    // Relations
    Toko        Toko           `json:"toko,omitempty" gorm:"foreignKey:IDUser"`
    Alamat      []Alamat       `json:"alamat,omitempty" gorm:"foreignKey:IDUser"`
}