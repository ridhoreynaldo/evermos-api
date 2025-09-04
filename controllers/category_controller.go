package controllers

import (
    "strconv"
    "github.com/gofiber/fiber/v2"
    "evermos-api/config"
    "evermos-api/models"
    "evermos-api/utils"
)

type CategoryController struct{}

type CategoryRequest struct {
    NamaCategory string `json:"nama_category" validate:"required"`
}

func (cc *CategoryController) CreateCategory(c *fiber.Ctx) error {
    var req CategoryRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    validate := utils.NewValidator()
    if err := validate.Struct(req); err != nil {
        return utils.ErrorResponse(c, 400, "Validation failed: "+err.Error())
    }

    // Check if category already exists
    var existingCategory models.Category
    if err := config.DB.Where("nama_category = ?", req.NamaCategory).First(&existingCategory).Error; err == nil {
        return utils.ErrorResponse(c, 400, "Category already exists")
    }

    category := models.Category{
        NamaCategory: req.NamaCategory,
    }

    if err := config.DB.Create(&category).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to create category")
    }

    return utils.SuccessResponse(c, 201, "Category created successfully", category)
}

func (cc *CategoryController) GetAllCategories(c *fiber.Ctx) error {
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    skip, _ := strconv.Atoi(c.Query("skip", "0"))
    nama := c.Query("nama")

    query := config.DB.Model(&models.Category{})

    if nama != "" {
        query = query.Where("nama_category LIKE ?", "%"+nama+"%")
    }

    var total int64
    query.Count(&total)

    var categories []models.Category
    if err := query.Limit(limit).Offset(skip).Find(&categories).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to fetch categories")
    }

    meta := utils.CalculateMeta(limit, skip, total)
    return utils.SuccessResponseWithMeta(c, 200, "Categories retrieved successfully", categories, meta)
}

func (cc *CategoryController) GetCategoryByID(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid category ID")
    }

    var category models.Category
    if err := config.DB.Preload("Products").First(&category, id).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Category not found")
    }

    return utils.SuccessResponse(c, 200, "Category retrieved successfully", category)
}

func (cc *CategoryController) UpdateCategory(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid category ID")
    }

    var req CategoryRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    validate := utils.NewValidator()
    if err := validate.Struct(req); err != nil {
        return utils.ErrorResponse(c, 400, "Validation failed: "+err.Error())
    }

    var category models.Category
    if err := config.DB.First(&category, id).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Category not found")
    }

    // Check if category name already exists (excluding current category)
    var existingCategory models.Category
    if err := config.DB.Where("nama_category = ? AND id != ?", req.NamaCategory, id).First(&existingCategory).Error; err == nil {
        return utils.ErrorResponse(c, 400, "Category name already exists")
    }

    category.NamaCategory = req.NamaCategory

    if err := config.DB.Save(&category).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to update category")
    }

    return utils.SuccessResponse(c, 200, "Category updated successfully", category)
}

func (cc *CategoryController) DeleteCategory(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid category ID")
    }

    var category models.Category
    if err := config.DB.First(&category, id).Error; err != nil {
        return utils.ErrorResponse(c, 404, "Category not found")
    }

    // Check if category has products
    var productCount int64
    config.DB.Model(&models.Product{}).Where("category_id = ?", id).Count(&productCount)
    if productCount > 0 {
        return utils.ErrorResponse(c, 400, "Cannot delete category with existing products")
    }

    if err := config.DB.Delete(&category).Error; err != nil {
        return utils.ErrorResponse(c, 500, "Failed to delete category")
    }

    return utils.SuccessResponse(c, 200, "Category deleted successfully", nil)
}
