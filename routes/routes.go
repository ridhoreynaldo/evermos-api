package routes

import (
    "github.com/gofiber/fiber/v2"
    "evermos-api/controllers"
    "evermos-api/middleware"
)

func SetupRoutes(app *fiber.App) {
    // Initialize controllers
    authController := &controllers.AuthController{}
    userController := &controllers.UserController{}
    tokoController := &controllers.TokoController{}
    alamatController := &controllers.AlamatController{}
    categoryController := &controllers.CategoryController{}
    productController := &controllers.ProductController{}
    transactionController := &controllers.TransactionController{}
    uploadController := &controllers.UploadController{}

    // Public routes
    api := app.Group("/api/v1")

    // Auth routes
    auth := api.Group("/auth")
    auth.Post("/register", authController.Register)
    auth.Post("/login", authController.Login)

    // Upload routes (public for now, you might want to add auth)
    upload := api.Group("/upload")
    upload.Post("/", uploadController.UploadFile)

    // Serve static files
    app.Static("/uploads", "./uploads")

    // Protected routes
    protected := api.Group("/", middleware.AuthMiddleware)

    // User routes
    user := protected.Group("/user")
    user.Get("/profile", userController.GetMyProfile)
    user.Put("/profile", userController.UpdateMyProfile)

    // Toko routes
    toko := api.Group("/toko")
    toko.Get("/", tokoController.GetAllToko)
    toko.Get("/:id", tokoController.GetTokoByID)
    
    protectedToko := protected.Group("/toko")
    protectedToko.Get("/my", tokoController.GetMyToko)
    protectedToko.Put("/my", tokoController.UpdateMyToko)

    // Alamat routes
    alamat := protected.Group("/alamat")
    alamat.Post("/", alamatController.CreateAlamat)
    alamat.Get("/", alamatController.GetMyAlamat)
    alamat.Get("/:id", alamatController.GetAlamatByID)
    alamat.Put("/:id", alamatController.UpdateAlamat)
    alamat.Delete("/:id", alamatController.DeleteAlamat)

    // Category routes (public read, admin write)
    category := api.Group("/category")
    category.Get("/", categoryController.GetAllCategories)
    category.Get("/:id", categoryController.GetCategoryByID)

    // Admin only category routes
    adminCategory := protected.Group("/category", middleware.AdminMiddleware)
    adminCategory.Post("/", categoryController.CreateCategory)
    adminCategory.Put("/:id", categoryController.UpdateCategory)
    adminCategory.Delete("/:id", categoryController.DeleteCategory)

    // Product routes
    product := api.Group("/product")
    product.Get("/", productController.GetAllProducts)
    product.Get("/:id", productController.GetProductByID)

    protectedProduct := protected.Group("/product")
    protectedProduct.Post("/", productController.CreateProduct)
    protectedProduct.Get("/my", productController.GetMyProducts)
    protectedProduct.Put("/:id", productController.UpdateProduct)
    protectedProduct.Delete("/:id", productController.DeleteProduct)

    // Transaction routes
    trx := protected.Group("/trx")
    trx.Post("/", transactionController.CreateTransaction)
    trx.Get("/", transactionController.GetMyTransactions)
    trx.Get("/:id", transactionController.GetTransactionByID)

    // Health check
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "OK",
            "message": "Evermos API is running",
        })
    })
}
