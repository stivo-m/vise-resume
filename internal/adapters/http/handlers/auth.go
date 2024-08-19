package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/stivo-m/vise-resume/internal/adapters/middleware"
	"github.com/stivo-m/vise-resume/internal/core/dto"
	"github.com/stivo-m/vise-resume/internal/core/ports"
	"github.com/stivo-m/vise-resume/internal/core/utils"
)

type AuthHandler struct {
	userService ports.UserService
}

func NewAuthHandler(
	userService ports.UserService,
) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (h AuthHandler) RegisterRoutes(router fiber.Router) {
	authRouter := router.Group("/auth")
	authRouter.Post(
		"/register",
		middleware.ValidationMiddleware(&dto.RegisterDto{}),
		h.HandleRegistration,
	)

	authRouter.Post(
		"/login",
		middleware.ValidationMiddleware(&dto.LoginDto{}),
		h.HandleLogin,
	)
}

func (h *AuthHandler) HandleRegistration(c *fiber.Ctx) error {
	var body dto.RegisterDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.userService.RegisterUser(context.Background(), body)
	if err != nil {
		data := utils.FormatApiResponse(
			"Registration failed",
			fiber.Map{"error": err.Error()},
		)
		return c.Status(fiber.StatusForbidden).JSON(data)
	}

	data := utils.FormatApiResponse(
		"User was registered successfully",
		res,
	)
	return c.Status(fiber.StatusCreated).JSON(data)
}

func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	var body dto.LoginDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.userService.LoginUser(context.Background(), body)
	if err != nil {
		data := utils.FormatApiResponse(
			"Login failed",
			fiber.Map{"error": err.Error()},
		)
		return c.Status(fiber.StatusForbidden).JSON(data)
	}

	data := utils.FormatApiResponse(
		"User was logged in successfully",
		res,
	)
	return c.Status(fiber.StatusOK).JSON(data)
}
