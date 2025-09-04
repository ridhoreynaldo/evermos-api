package controllers

import (
    "strconv"
    "github.com/gofiber/fiber/v2"
    "evermos-api/config"
    "evermos-api/models"
    "evermos-api/utils"
)

type TokoController struct{}

type TokoRequest struct {
    NamaToko string `json:"nama_toko" validate:"required"`
    UrlFoto  string `json:"url_foto"`
}

func (tc *TokoController) GetMyToko(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)

    var toko models.Toko
    if err := config.DB.Where("id_user = ?", userID).Preload("User").First(&toko).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Store not found")
    }

    return utils.SuccessResponse(c, 200, "Store retrieved successfully", toko)
}

func (tc *TokoController) UpdateMyToko(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)

    var req TokoRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    validate := utils.NewValidator()
    if err := validate.Struct(req); err != nil {
        return utils.ErrorResponse(c, 400, "Validation failed: "+err.Error())
    }

    var toko models.Toko
    if err := config.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Store not found")
    }

    toko.NamaToko = req.NamaToko
    toko.UrlFoto = req.UrlFoto

    if err := config.DB.Save(&toko).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to update store")
    }

    return utils.SuccessResponse(c, 200, "Store updated successfully", toko)
}

func (tc *TokoController) GetTokoByID(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid store ID")
    }

    var toko models.Toko
    if err := config.DB.Preload("User").Preload("Products").First(&toko, id).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Store not found")
    }

    return utils.SuccessResponse(c, 200, "Store retrieved successfully", toko)
}

func (tc *TokoController) GetAllToko(c *fiber.Ctx) error {
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    skip, _ := strconv.Atoi(c.Query("skip", "0"))
    nama := c.Query("nama")

    query := config.DB.Preload("User")

    if nama != "" {
        query = query.Where("nama_toko LIKE ?", "%"+nama+"%")
    }

    var total int64
    query.Model(&models.Toko{}).Count(&total)

    var tokos []models.Toko
    if err := query.Limit(limit).Offset(skip).Find(&tokos).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to fetch stores")
    }

    meta := utils.CalculateMeta(limit, skip, total)
    return utils.SuccessResponseWithMeta(c, 200, "Stores retrieved successfully", tokos, meta)
}
