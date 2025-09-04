package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
    Status  int         `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Meta    interface{} `json:"meta,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
    return c.Status(status).JSON(Response{
        Status:  status,
        Message: message,
        Data:    data,
    })
}

func SuccessResponseWithMeta(c *fiber.Ctx, status int, message string, data interface{}, meta interface{}) error {
    return c.Status(status).JSON(Response{
        Status:  status,
        Message: message,
        Data:    data,
        Meta:    meta,
    })
}

func ErrorResponse(c *fiber.Ctx, status int, message string) error {
    return c.Status(status).JSON(Response{
        Status:  status,
        Message: message,
    })
}

type Meta struct {
    Limit      int   `json:"limit"`
    Skip       int   `json:"skip"`
    Total      int64 `json:"total"`
    TotalPages int   `json:"total_pages"`
}

func CalculateMeta(limit, skip int, total int64) Meta {
    totalPages := int(total) / limit
    if int(total)%limit > 0 {
        totalPages++
    }
    
    return Meta{
        Limit:      limit,
        Skip:       skip,
        Total:      total,
        TotalPages: totalPages,
    }
}
