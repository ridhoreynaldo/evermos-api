package controllers

import (
    "fmt"
    "strings"
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "evermos-api/config"
    "evermos-api/models"
    "evermos-api/utils"
)

type AuthController struct{}

type RegisterRequest struct {
    Nama         string    `json:"nama" validate:"required"`
    KataSandi    string    `json:"kata_sandi" validate:"required,min=6"`
    NoTelp       string    `json:"no_telp" validate:"required"`
    TanggalLahir time.Time `json:"tanggal_lahir" validate:"required"`
    JenisKelamin string    `json:"jenis_kelamin" validate:"required,oneof=L P"`
    Tentang      string    `json:"tentang"`
    Pekerjaan    string    `json:"pekerjaan"`
    Email        string    `json:"email" validate:"required,email"`
    IDProvinsi   string    `json:"id_provinsi" validate:"required"`
    IDKota       string    `json:"id_kota" validate:"required"`
}

type LoginRequest struct {
    Email     string `json:"email" validate:"required,email"`
    KataSandi string `json:"kata_sandi" validate:"required"`
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    // Validate request
    validate := utils.NewValidator()
    if err := validate.Struct(req); err != nil {
        return utils.ErrorResponse(c, 400, "Validation failed: "+err.Error())
    }

    // Check if email or phone already exists
    var existingUser models.User
    if err := config.DB.Where("email = ? OR no_telp = ?", req.Email, req.NoTelp).First(&existingUser).Error; err == nil {
        return utils.ErrorResponse(c, 400, "Email or phone number already registered")
    }

    // Hash password
    hashedPassword, err := utils.HashPassword(req.KataSandi)
    if err != nil {
        return utils.ErrorResponse(c, 500, "Failed to hash password")
    }

    // Create user
    user := models.User{
        Nama:         req.Nama,
        KataSandi:    hashedPassword,
        NoTelp:       req.NoTelp,
        TanggalLahir: req.TanggalLahir,
        JenisKelamin: req.JenisKelamin,
        Tentang:      req.Tentang,
        Pekerjaan:    req.Pekerjaan,
        Email:        req.Email,
        IDProvinsi:   req.IDProvinsi,
        IDKota:       req.IDKota,
        IsAdmin:      false,
    }

    tx := config.DB.Begin()
    
    if err := tx.Create(&user).Error; err != nil {
        tx.Rollback()
        return utils.ErrorResponse(c, 500, "Failed to create user")
    }

    // Create toko automatically
    toko := models.Toko{
        NamaToko: user.Nama + "'s Store",
        IDUser:   user.ID,
    }

    if err := tx.Create(&toko).Error; err != nil {
        tx.Rollback()
        return utils.ErrorResponse(c, 500, "Failed to create store")
    }

    tx.Commit()

    // Generate JWT token
    token, err := utils.GenerateJWT(user.ID, user.Email, user.IsAdmin)
    if err != nil {
        return utils.ErrorResponse(c, 500, "Failed to generate token")
    }

    // Remove password from response
    user.KataSandi = ""

    return utils.SuccessResponse(c, 201, "Registration successful", fiber.Map{
        "user":  user,
        "token": token,
    })
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    // Validate request
    validate := utils.NewValidator()
    if err := validate.Struct(req); err != nil {
        return utils.ErrorResponse(c, 400, "Validation failed: "+err.Error())
    }

    // Find user
    var user models.User
    if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        return utils.ErrorResponse(c, 401, "Invalid email or password")
    }

    // Check password
    if !utils.CheckPasswordHash(req.KataSandi, user.KataSandi) {
        return utils.ErrorResponse(c, 401, "Invalid email or password")
    }

    // Generate JWT token
    token, err := utils.GenerateJWT(user.ID, user.Email, user.IsAdmin)
    if err != nil {
        return utils.ErrorResponse(c, 500, "Failed to generate token")
    }

    // Remove password from response
    user.KataSandi = ""

    return utils.SuccessResponse(c, 200, "Login successful", fiber.Map{
        "user":  user,
        "token": token,
    })
}
