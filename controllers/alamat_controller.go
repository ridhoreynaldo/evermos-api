package controllers

import (
    "strconv"
    "github.com/gofiber/fiber/v2"
    "evermos-api/config"
    "evermos-api/models"
    "evermos-api/utils"
)

type AlamatController struct{}

type AlamatRequest struct {
    JudulAlamat  string `json:"judul_alamat" validate:"required"`
    NamaPenerima string `json:"nama_penerima" validate:"required"`
    NoTelp       string `json:"no_telp" validate:"required"`
    DetailAlamat string `json:"detail_alamat" validate:"required"`
}

func (ac *AlamatController) CreateAlamat(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)

    var req AlamatRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    validate := utils.NewValidator()
    if err := validate.Struct(req); err != nil {
        return utils.ErrorResponse(c, 400, "Validation failed: "+err.Error())
    }

    alamat := models.Alamat{
        JudulAlamat:  req.JudulAlamat,
        NamaPenerima: req.NamaPenerima,
        NoTelp:       req.NoTelp,
        DetailAlamat: req.DetailAlamat,
        IDUser:       userID,
    }

    if err := config.DB.Create(&alamat).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to create address")
    }

    return utils.SuccessResponse(c, 201, "Address created successfully", alamat)
}

func (ac *AlamatController) GetMyAlamat(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)

    var alamats []models.Alamat
    if err := config.DB.Where("id_user = ?", userID).Find(&alamats).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to fetch addresses")
    }

    return utils.SuccessResponse(c, 200, "Addresses retrieved successfully", alamats)
}

func (ac *AlamatController) GetAlamatByID(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid address ID")
    }

    var alamat models.Alamat
    if err := config.DB.Where("id = ? AND id_user = ?", id, userID).First(&alamat).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Address not found")
    }

    return utils.SuccessResponse(c, 200, "Address retrieved successfully", alamat)
}

func (ac *AlamatController) UpdateAlamat(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid address ID")
    }

    var req AlamatRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    validate := utils.NewValidator()
    if err := validate.Struct(req); err != nil {
        return utils.ErrorResponse(c, 400, "Validation failed: "+err.Error())
    }

    var alamat models.Alamat
    if err := config.DB.Where("id = ? AND id_user = ?", id, userID).First(&alamat).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Address not found")
    }

    alamat.JudulAlamat = req.JudulAlamat
    alamat.NamaPenerima = req.NamaPenerima
    alamat.NoTelp = req.NoTelp
    alamat.DetailAlamat = req.DetailAlamat

    if err := config.DB.Save(&alamat).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to update address")
    }

    return utils.SuccessResponse(c, 200, "Address updated successfully", alamat)
}

func (ac *AlamatController) DeleteAlamat(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid address ID")
    }

    var alamat models.Alamat
    if err := config.DB.Where("id = ? AND id_user = ?", id, userID).First(&alamat).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Address not found")
    }

    if err := config.DB.Delete(&alamat).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to delete address")
    }

    return utils.SuccessResponse(c, 200, "Address deleted successfully", nil)
}
