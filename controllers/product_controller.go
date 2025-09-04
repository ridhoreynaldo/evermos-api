package controllers

import (
    "strconv"
    "github.com/gofiber/fiber/v2"
    "evermos-api/config"
    "evermos-api/models"
    "evermos-api/utils"
)

type ProductController struct{}

type ProductRequest struct {
    NamaProduk     string   `json:"nama_produk" validate:"required"`
    HargaReseller  int      `json:"harga_reseller" validate:"required,min=0"`
    HargaKonsumen  int      `json:"harga_konsumen" validate:"required,min=0"`
    Stok           int      `json:"stok" validate:"required,min=0"`
    Deskripsi      string   `json:"deskripsi"`
    URLFoto        []string `json:"url_foto"`
    CategoryID     uint     `json:"category_id" validate:"required"`
}

func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)

    var req ProductRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    validate := utils.NewValidator()
    if err := validate.Struct(req); err != nil {
        return utils.ErrorResponse(c, 400, "Validation failed: "+err.Error())
    }

    // Get user's toko
    var toko models.Toko
    if err := config.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Store not found")
    }

    // Check if category exists
    var category models.Category
    if err := config.DB.First(&category, req.CategoryID).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Category not found")
    }

    // Generate slug
    slug := utils.GenerateSlug(req.NamaProduk)

    product := models.Product{
        NamaProduk:    req.NamaProduk,
        Slug:          slug,
        HargaReseller: req.HargaReseller,
        HargaKonsumen: req.HargaKonsumen,
        Stok:          req.Stok,
        Deskripsi:     req.Deskripsi,
        URLFoto:       req.URLFoto,
        IDToko:        toko.ID,
        CategoryID:    req.CategoryID,
    }

    if err := config.DB.Create(&product).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to create product")
    }

    // Preload relations
    config.DB.Preload("Toko").Preload("Category").First(&product, product.ID)

    return utils.SuccessResponse(c, 201, "Product created successfully", product)
}

func (pc *ProductController) GetMyProducts(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    skip, _ := strconv.Atoi(c.Query("skip", "0"))
    nama := c.Query("nama")
    categoryID := c.Query("category_id")

    // Get user's toko
    var toko models.Toko
    if err := config.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Store not found")
    }

    query := config.DB.Preload("Toko").Preload("Category").Where("id_toko = ?", toko.ID)

    if nama != "" {
        query = query.Where("nama_produk LIKE ?", "%"+nama+"%")
    }

    if categoryID != "" {
        query = query.Where("category_id = ?", categoryID)
    }

    var total int64
    query.Model(&models.Product{}).Count(&total)

    var products []models.Product
    if err := query.Limit(limit).Offset(skip).Find(&products).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to fetch products")
    }

    meta := utils.CalculateMeta(limit, skip, total)
    return utils.SuccessResponseWithMeta(c, 200, "Products retrieved successfully", products, meta)
}

func (pc *ProductController) GetAllProducts(c *fiber.Ctx) error {
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    skip, _ := strconv.Atoi(c.Query("skip", "0"))
    nama := c.Query("nama")
    categoryID := c.Query("category_id")
    tokoID := c.Query("toko_id")

    query := config.DB.Preload("Toko").Preload("Category")

    if nama != "" {
        query = query.Where("nama_produk LIKE ?", "%"+nama+"%")
    }

    if categoryID != "" {
        query = query.Where("category_id = ?", categoryID)
    }

    if tokoID != "" {
        query = query.Where("id_toko = ?", tokoID)
    }

    var total int64
    query.Model(&models.Product{}).Count(&total)

    var products []models.Product
    if err := query.Limit(limit).Offset(skip).Find(&products).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to fetch products")
    }

    meta := utils.CalculateMeta(limit, skip, total)
    return utils.SuccessResponseWithMeta(c, 200, "Products retrieved successfully", products, meta)
}

func (pc *ProductController) GetProductByID(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid product ID")
    }

    var product models.Product
    if err := config.DB.Preload("Toko").Preload("Category").First(&product, id).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Product not found")
    }

    return utils.SuccessResponse(c, 200, "Product retrieved successfully", product)
}

func (pc *ProductController) UpdateProduct(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid product ID")
    }

    var req ProductRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    validate := utils.NewValidator()
    if err := validate.Struct(req); err != nil {
        return utils.ErrorResponse(c, 400, "Validation failed: "+err.Error())
    }

    // Get user's toko
    var toko models.Toko
    if err := config.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Store not found")
    }

    var product models.Product
    if err := config.DB.Where("id = ? AND id_toko = ?", id, toko.ID).First(&product).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Product not found")
    }

    // Check if category exists
    var category models.Category
    if err := config.DB.First(&category, req.CategoryID).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Category not found")
    }

    product.NamaProduk = req.NamaProduk
    product.HargaReseller = req.HargaReseller
    product.HargaKonsumen = req.HargaKonsumen
    product.Stok = req.Stok
    product.Deskripsi = req.Deskripsi
    product.URLFoto = req.URLFoto
    product.CategoryID = req.CategoryID

    // Update slug if name changed
    if req.NamaProduk != product.NamaProduk {
        product.Slug = utils.GenerateSlug(req.NamaProduk)
    }

    if err := config.DB.Save(&product).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to update product")
    }

    // Preload relations
    config.DB.Preload("Toko").Preload("Category").First(&product, product.ID)

    return utils.SuccessResponse(c, 200, "Product updated successfully", product)
}

func (pc *ProductController) DeleteProduct(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid product ID")
    }

    // Get user's toko
    var toko models.Toko
    if err := config.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Store not found")
    }

    var product models.Product
    if err := config.DB.Where("id = ? AND id_toko = ?", id, toko.ID).First(&product).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Product not found")
    }

    if err := config.DB.Delete(&product).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to delete product")
    }

    return utils.SuccessResponse(c, 200, "Product deleted successfully", nil)
}