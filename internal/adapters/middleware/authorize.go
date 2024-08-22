package middleware

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/stivo-m/vise-resume/internal/core/dto"
	"github.com/stivo-m/vise-resume/internal/core/ports"
	"github.com/stivo-m/vise-resume/internal/core/utils"
)

func AuthMiddleware(tokenPort ports.TokenService, userPort ports.UserPort) fiber.Handler {
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

		record, err := userPort.FindToken(context.Background(), dto.ManageTokenDto{
			ID:          userId,
			AccessToken: tokenString,
		})

		if err != nil {
			res := utils.FormatApiResponse(
				"authentication failed",
				err,
			)
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}

		if record == nil || record.DeletedAt.Valid {
			res := utils.FormatApiResponse(
				"authentication failed",
				nil,
			)
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}

		utils.TextLogger.Info("----------------------------------------------------------------")
		utils.TextLogger.Info(fmt.Sprintf("authenticated user: %s", userId))
		utils.TextLogger.Info("----------------------------------------------------------------")

		c.Locals(utils.USER_ID_KEY, userId)
		c.Locals(utils.ACCESS_TOKEN_KEY, tokenString)
		return c.Next()
	}
}
