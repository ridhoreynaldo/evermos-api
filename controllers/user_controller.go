package controllers

import (
    "strconv"
    "github.com/gofiber/fiber/v2"
    "evermos-api/config"
    "evermos-api/models"
    "evermos-api/utils"
)

type UserController struct{}

func (uc *UserController) GetMyProfile(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)

    var user models.User
    if err := config.DB.Preload("Toko").First(&user, userID).Error; err != nil {
        return utils.ErrorResponse(c, 404, "User not found")
    }

    user.KataSandi = ""
    return utils.SuccessResponse(c, 200, "Profile retrieved successfully", user)
}

func (uc *UserController) UpdateMyProfile(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)

    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    var user models.User
    if err := config.DB.First(&user, userID).Error; err != nil {
        return utils.ErrorResponse(c, 404, "User not found")
    }

    // Check if email or phone is taken by another user
    var existingUser models.User
    if err := config.DB.Where("(email = ? OR no_telp = ?) AND id != ?", req.Email, req.NoTelp, userID).First(&existingUser).Error; err == nil {
        return utils.ErrorResponse(c, 400, "Email or phone number already taken")
    }

    // Update user fields
    user.Nama = req.Nama
    user.NoTelp = req.NoTelp
    user.TanggalLahir = req.TanggalLahir
    user.JenisKelamin = req.JenisKelamin
    user.Tentang = req.Tentang
    user.Pekerjaan = req.Pekerjaan
    user.Email = req.Email
    user.IDProvinsi = req.IDProvinsi
    user.IDKota = req.IDKota

    // Update password if provided
    if req.KataSandi != "" {
        hashedPassword, err := utils.HashPassword(req.KataSandi)
        if err != nil {
            return utils.ErrorResponse(c, 500, "Failed to hash password")
        }
        user.KataSandi = hashedPassword
    }

    if err := config.DB.Save(&user).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to update profile")
    }

    user.KataSandi = ""
    return utils.SuccessResponse(c, 200, "Profile updated successfully", user)
}
