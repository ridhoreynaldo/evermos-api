package config

import (
    "fmt"
    "log"
    "os"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "evermos-api/models"
)

var DB *gorm.DB

func ConnectDB() {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        getEnv("DB_USER", "root"),
        getEnv("DB_PASSWORD", ""),
        getEnv("DB_HOST", "localhost"),
        getEnv("DB_PORT", "3306"),
        getEnv("DB_NAME", "evermos_db"),
    )

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto Migrate
    db.AutoMigrate(
        &models.User{},
        &models.Toko{},
        &models.Alamat{},
        &models.Category{},
        &models.Product{},
        &models.Trx{},
        &models.DetailTrx{},
        &models.LogProduct{},
    )

    DB = db
    fmt.Println("Database connected successfully")
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
