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
	tokenPort   ports.TokenService
}

func NewAuthHandler(
	userService ports.UserService,
	tokenPort ports.TokenService,
) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		tokenPort:   tokenPort,
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

	authRouter.Post(
		"/verify-email",
		middleware.ValidationMiddleware(&dto.VerificationDto{}),
		h.HandleEmailVerification,
	)

	authRouter.Post(
		"/forgot-password",
		middleware.ValidationMiddleware(&dto.EmailDto{}),
		h.HandleForgotPassword,
	)

	authRouter.Post(
		"/reset-password",
		middleware.ValidationMiddleware(&dto.ResetPasswordDto{}),
		h.HandleResetPassword,
	)

	authRouter.Post(
		"/logout",
		middleware.AuthMiddleware(h.tokenPort),
		h.HandleLogout,
	)
}

// Handles the process of registering a user
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

// Handles the process of authenticating a user
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

// Handles the process of verifying a user's account
func (h *AuthHandler) HandleEmailVerification(c *fiber.Ctx) error {
	var body dto.VerificationDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.userService.VerifyEmailAddress(context.Background(), body)
	if err != nil {
		data := utils.FormatApiResponse(
			"Verification failed",
			fiber.Map{"error": err.Error()},
		)
		return c.Status(fiber.StatusForbidden).JSON(data)
	}

	data := utils.FormatApiResponse("Email verification successful", nil)
	return c.Status(fiber.StatusOK).JSON(data)
}

// Handles the process of forgetting a user's password
func (h *AuthHandler) HandleForgotPassword(c *fiber.Ctx) error {
	var body dto.EmailDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.userService.ForgetPassword(context.Background(), body)
	if err != nil {
		data := utils.FormatApiResponse(
			"Password reset failed",
			fiber.Map{"error": err.Error()},
		)
		return c.Status(fiber.StatusForbidden).JSON(data)
	}

	data := utils.FormatApiResponse(
		"Instructions have been sent to the registered email address on how to reset your password",
		nil,
	)
	return c.Status(fiber.StatusOK).JSON(data)
}

// Handles the process of resetting a user's account
func (h *AuthHandler) HandleResetPassword(c *fiber.Ctx) error {
	var body dto.ResetPasswordDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.userService.ResetPassword(context.Background(), body)
	if err != nil {
		data := utils.FormatApiResponse(
			"Password reset failed",
			fiber.Map{"error": err.Error()},
		)
		return c.Status(fiber.StatusForbidden).JSON(data)
	}

	data := utils.FormatApiResponse(
		"Instructions have been sent to the registered email address on how to reset your password",
		nil,
	)
	return c.Status(fiber.StatusOK).JSON(data)
}

// Handles the process of resetting a user's account
func (h *AuthHandler) HandleLogout(c *fiber.Ctx) error {
	accessToken := c.Locals(utils.ACCESS_TOKEN_KEY).(string)
	err := h.userService.LogoutUser(context.Background(), accessToken)
	if err != nil {
		data := utils.FormatApiResponse(
			"Password reset failed",
			fiber.Map{"error": err.Error()},
		)
		return c.Status(fiber.StatusForbidden).JSON(data)
	}

	data := utils.FormatApiResponse(
		"Instructions have been sent to the registered email address on how to reset your password",
		nil,
	)
	return c.Status(fiber.StatusOK).JSON(data)
}
