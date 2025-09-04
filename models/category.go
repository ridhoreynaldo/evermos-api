package models

import (
    "time"
)

type Category struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    NamaCategory string    `json:"nama_category" gorm:"not null;unique" validate:"required"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    
    // Relations
    Products     []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}