package controllers

import (
    "fmt"
    "mime/multipart"
    "path/filepath"
    "strings"
    "time"
    "github.com/gofiber/fiber/v2"
    "evermos-api/utils"
)

type UploadController struct{}

func (uc *UploadController) UploadFile(c *fiber.Ctx) error {
    // Parse multipart form
    form, err := c.MultipartForm()
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid multipart form")
    }

    files := form.File["file"]
    if len(files) == 0 {
        return utils.ErrorResponse(c, 400, "No file uploaded")
    }

    var uploadedFiles []string

    for _, file := range files {
        // Validate file type
        if !isValidFileType(file) {
            return utils.ErrorResponse(c, 400, "Invalid file type. Only images are allowed")
        }

        // Validate file size (max 5MB)
        if file.Size > 5*1024*1024 {
            return utils.ErrorResponse(c, 400, "File size too large. Maximum 5MB allowed")
        }

        // Generate unique filename
        filename := generateUniqueFilename(file.Filename)
        uploadPath := fmt.Sprintf("./uploads/%s", filename)

        // Save file
        if err := c.SaveFile(file, uploadPath); err != nil {
            return utils.ErrorResponse(c, 500, "Failed to save file")
        }

        // Add to response (you might want to use a proper URL with domain)
        fileURL := fmt.Sprintf("/uploads/%s", filename)
        uploadedFiles = append(uploadedFiles, fileURL)
    }

    return utils.SuccessResponse(c, 200, "Files uploaded successfully", fiber.Map{
        "files": uploadedFiles,
    })
}

func isValidFileType(file *multipart.FileHeader) bool {
    allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
    ext := strings.ToLower(filepath.Ext(file.Filename))
    
    for _, allowed := range allowedExtensions {
        if ext == allowed {
            return true
        }
    }
    
    return false
}

func generateUniqueFilename(originalFilename string) string {
    ext := filepath.Ext(originalFilename)
    timestamp := time.Now().Unix()
    return fmt.Sprintf("%d_%s%s", timestamp, utils.GenerateRandomString(8), ext)
}
