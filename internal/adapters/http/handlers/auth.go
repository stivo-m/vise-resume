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
	userPort    ports.UserPort
	tokenPort   ports.TokenService
}

func NewAuthHandler(
	userService ports.UserService,
	userPort ports.UserPort,
	tokenPort ports.TokenService,
) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		userPort:    userPort,
		tokenPort:   tokenPort,
	}
}

func (h AuthHandler) RegisterRoutes(router fiber.Router) {
	authRouter := router.Group("/auth")
	authRouter.Post(
		"/register",
		middleware.ValidationMiddleware(&dto.RegisterDto{}),
		h.handleRegistration,
	)

	authRouter.Post(
		"/login",
		middleware.ValidationMiddleware(&dto.LoginDto{}),
		h.handleLogin,
	)

	authRouter.Post(
		"/verify-email",
		middleware.ValidationMiddleware(&dto.VerificationDto{}),
		h.handleEmailVerification,
	)

	authRouter.Post(
		"/forgot-password",
		middleware.ValidationMiddleware(&dto.EmailDto{}),
		h.handleForgotPassword,
	)

	authRouter.Post(
		"/reset-password",
		middleware.ValidationMiddleware(&dto.ResetPasswordDto{}),
		h.handleResetPassword,
	)

	authRouter.Post(
		"/logout",
		middleware.AuthMiddleware(h.tokenPort, h.userPort),
		h.handleLogout,
	)

	authRouter.Get(
		"/profile",
		middleware.AuthMiddleware(h.tokenPort, h.userPort),
		h.handleShowProfile,
	)

	authRouter.Patch(
		"/profile",
		middleware.ValidationMiddleware(&dto.UpdateUserDto{}),
		middleware.AuthMiddleware(h.tokenPort, h.userPort),
		h.handleUpdateUserInfo,
	)
}

// Handles the process of registering a user
func (h *AuthHandler) handleRegistration(c *fiber.Ctx) error {
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
func (h *AuthHandler) handleLogin(c *fiber.Ctx) error {
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
func (h *AuthHandler) handleEmailVerification(c *fiber.Ctx) error {
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
func (h *AuthHandler) handleForgotPassword(c *fiber.Ctx) error {
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
func (h *AuthHandler) handleResetPassword(c *fiber.Ctx) error {
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
		"Password reset successful. Your password has been updated and will be used to authenticate in the future.",
		nil,
	)
	return c.Status(fiber.StatusOK).JSON(data)
}

// Handles the process of resetting a user's account
func (h *AuthHandler) handleLogout(c *fiber.Ctx) error {
	accessToken := c.Locals(utils.ACCESS_TOKEN_KEY).(string)
	err := h.userService.LogoutUser(context.Background(), accessToken)
	if err != nil {
		data := utils.FormatApiResponse(
			"Logout failed",
			fiber.Map{"error": err.Error()},
		)
		return c.Status(fiber.StatusForbidden).JSON(data)
	}

	data := utils.FormatApiResponse(
		"Logout was successful",
		nil,
	)
	return c.Status(fiber.StatusOK).JSON(data)
}

// Handles the process of showing a user's profile
func (h *AuthHandler) handleShowProfile(c *fiber.Ctx) error {
	userId := c.Locals(utils.USER_ID_KEY).(string)
	user, err := h.userService.ShowProfile(context.Background(), userId)
	if err != nil {
		data := utils.FormatApiResponse(
			"Unable to show profile",
			fiber.Map{"error": err.Error()},
		)
		return c.Status(fiber.StatusForbidden).JSON(data)
	}

	data := utils.FormatApiResponse(
		"User profile",
		user,
	)
	return c.Status(fiber.StatusOK).JSON(data)
}

// Handles the process of updating a user's name
func (h *AuthHandler) handleUpdateUserInfo(c *fiber.Ctx) error {
	var body dto.UpdateUserDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updates := map[string]interface{}{
		"full_name": body.FullName,
	}
	userId := c.Locals(utils.USER_ID_KEY).(string)

	err := h.userService.UpdateUser(context.Background(), userId, updates)
	if err != nil {
		data := utils.FormatApiResponse(
			"User update failed",
			fiber.Map{"error": err.Error()},
		)
		return c.Status(fiber.StatusForbidden).JSON(data)
	}

	data := utils.FormatApiResponse(
		"User updated successfully",
		nil,
	)
	return c.Status(fiber.StatusOK).JSON(data)
}
