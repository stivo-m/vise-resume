package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stivo-m/vise-resume/internal/core/ports"
	"github.com/stivo-m/vise-resume/internal/core/utils"
)

func AuthMiddleware(tokenPort ports.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for authorization header and token
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			res := utils.FormatApiResponse(
				"unauthorized",
				nil,
			)
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}

		// Bearer token format: "Bearer <token>"
		tokenString = tokenString[len("Bearer "):]
		userId, err := tokenPort.VerifyToken(tokenString)

		if err != nil {
			res := utils.FormatApiResponse(
				"authentication failed",
				err,
			)
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}

		c.Locals(utils.USER_ID_KEY, userId)
		c.Locals(utils.ACCESS_TOKEN_KEY, tokenString)
		return c.Next()
	}
}
