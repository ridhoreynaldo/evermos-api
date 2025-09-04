package utils

import (
    "regexp"
    "strings"
    "github.com/google/uuid"
)

func GenerateSlug(text string) string {
    // Convert to lowercase
    text = strings.ToLower(text)
    
    // Replace spaces with hyphens
    text = strings.ReplaceAll(text, " ", "-")
    
    // Remove special characters
    reg := regexp.MustCompile(`[^a-z0-9-]`)
    text = reg.ReplaceAllString(text, "")
    
    // Remove multiple consecutive hyphens
    reg = regexp.MustCompile(`-+`)
    text = reg.ReplaceAllString(text, "-")
    
    // Trim hyphens from start and end
    text = strings.Trim(text, "-")
    
    // Add UUID to make it unique
    if text == "" {
        text = "product"
    }
    
    return text + "-" + uuid.New().String()[:8]
}