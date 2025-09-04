package controllers

import (
    "fmt"
    "strconv"
    "time"
    "github.com/gofiber/fiber/v2"
    "evermos-api/config"
    "evermos-api/models"
    "evermos-api/utils"
)

type TransactionController struct{}

type TransactionRequest struct {
    AlamatKirim uint                      `json:"alamat_kirim" validate:"required"`
    MetodeBayar string                    `json:"metode_bayar" validate:"required"`
    DetailTrx   []DetailTransactionRequest `json:"detail_trx" validate:"required,dive"`
}

type DetailTransactionRequest struct {
    IDProduk  uint `json:"id_produk" validate:"required"`
    Kuantitas int  `json:"kuantitas" validate:"required,min=1"`
}

func (tc *TransactionController) CreateTransaction(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)

    var req TransactionRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    validate := utils.NewValidator()
    if err := validate.Struct(req); err != nil {
        return utils.ErrorResponse(c, 400, "Validation failed: "+err.Error())
    }

    // Check if alamat belongs to user
    var alamat models.Alamat
    if err := config.DB.Where("id = ? AND id_user = ?", req.AlamatKirim, userID).First(&alamat).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Address not found")
    }

    tx := config.DB.Begin()

    // Calculate total price and validate products
    var totalPrice int
    var detailTrxs []models.DetailTrx
    
    for _, detail := range req.DetailTrx {
        var product models.Product
        if err := tx.First(&product, detail.IDProduk).Error; err != nil {
            tx.Rollback()
            return utils.ErrorResponse(c, 404, "Product not found")
        }

        // Check stock
        if product.Stok < detail.Kuantitas {
            tx.Rollback()
            return utils.ErrorResponse(c, 400, fmt.Sprintf("Insufficient stock for product %s", product.NamaProduk))
        }

        // Create log product
        logProduct := models.LogProduct{
            IDProduk:      product.ID,
            NamaProduk:    product.NamaProduk,
            Slug:          product.Slug,
            HargaReseller: product.HargaReseller,
            HargaKonsumen: product.HargaKonsumen,
            Deskripsi:     product.Deskripsi,
            IDToko:        product.IDToko,
            CategoryID:    product.CategoryID,
        }

        if err := tx.Create(&logProduct).Error; err != nil {
            tx.Rollback()
            return utils.ErrorResponse(c, 500, "Failed to create product log")
        }

        // Calculate detail total
        detailTotal := product.HargaKonsumen * detail.Kuantitas
        totalPrice += detailTotal

        detailTrx := models.DetailTrx{
            IDLogProduk: logProduct.ID,
            IDToko:      product.IDToko,
            Kuantitas:   detail.Kuantitas,
            HargaTotal:  detailTotal,
        }

        detailTrxs = append(detailTrxs, detailTrx)

        // Update product stock
        product.Stok -= detail.Kuantitas
        if err := tx.Save(&product).Error; err != nil {
            tx.Rollback()
            return utils.ErrorResponse(c, 500, "Failed to update product stock")
        }
    }

    // Generate invoice code
    kodeInvoice := fmt.Sprintf("INV-%d-%d", userID, time.Now().Unix())

    // Create transaction
    trx := models.Trx{
        IDUser:      userID,
        AlamatKirim: req.AlamatKirim,
        HargaTotal:  totalPrice,
        KodeInvoice: kodeInvoice,
        MetodeBayar: req.MetodeBayar,
    }

    if err := tx.Create(&trx).Error; err != nil {
        tx.Rollback()
        return utils.ErrorResponse(c, 500, "Failed to create transaction")
    }

    // Create detail transactions
    for i := range detailTrxs {
        detailTrxs[i].IDTrx = trx.ID
        if err := tx.Create(&detailTrxs[i]).Error; err != nil {
            tx.Rollback()
            return utils.ErrorResponse(c, 500, "Failed to create transaction details")
        }
    }

    tx.Commit()

    // Preload relations
    config.DB.Preload("User").Preload("Alamat").Preload("DetailTrx").Preload("DetailTrx.LogProduct").Preload("DetailTrx.Toko").First(&trx, trx.ID)

    return utils.SuccessResponse(c, 201, "Transaction created successfully", trx)
}

func (tc *TransactionController) GetMyTransactions(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    skip, _ := strconv.Atoi(c.Query("skip", "0"))

    query := config.DB.Preload("User").Preload("Alamat").Preload("DetailTrx").Preload("DetailTrx.LogProduct").Preload("DetailTrx.Toko").Where("id_user = ?", userID)

    var total int64
    query.Model(&models.Trx{}).Count(&total)

    var transactions []models.Trx
    if err := query.Limit(limit).Offset(skip).Find(&transactions).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to fetch transactions")
    }

    meta := utils.CalculateMeta(limit, skip, total)
    return utils.SuccessResponseWithMeta(c, 200, "Transactions retrieved successfully", transactions, meta)
}

func (tc *TransactionController) GetTransactionByID(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid transaction ID")
    }

    var transaction models.Trx
    if err := config.DB.Preload("User").Preload("Alamat").Preload("DetailTrx").Preload("DetailTrx.LogProduct").Preload("DetailTrx.Toko").Where("id = ? AND id_user = ?", id, userID).First(&transaction).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Transaction not found")
    }

    return utils.SuccessResponse(c, 200, "Transaction retrieved successfully", transaction)
}
