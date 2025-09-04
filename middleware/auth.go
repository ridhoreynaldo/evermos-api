package middleware

import (
    "strings"
    "github.com/gofiber/fiber/v2"
    "evermos-api/utils"
)

func AuthMiddleware(c *fiber.Ctx) error {
    token := c.Get("Authorization")
    if token == "" {
        return utils.ErrorResponse(c, 401, "Authorization header required")
    }

    // Remove "Bearer " prefix
    if strings.HasPrefix(token, "Bearer ") {
        token = strings.TrimPrefix(token, "Bearer ")
    }

    claims, err := utils.ValidateJWT(token)
    if err != nil {
        return utils.ErrorResponse(c, 401, "Invalid token")
    }

    // Set user info in context
    c.Locals("user_id", claims.UserID)
    c.Locals("email", claims.Email)
    c.Locals("is_admin", claims.IsAdmin)

    return c.Next()
}

func AdminMiddleware(c *fiber.Ctx) error {
    isAdmin, ok := c.Locals("is_admin").(bool)
    if !ok || !isAdmin {
        return utils.ErrorResponse(c, 403, "Admin access required")
    }
    
    return c.Next()
}
