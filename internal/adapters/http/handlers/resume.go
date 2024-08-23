package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/stivo-m/vise-resume/internal/adapters/middleware"
	"github.com/stivo-m/vise-resume/internal/core/dto"
	"github.com/stivo-m/vise-resume/internal/core/ports"
	"github.com/stivo-m/vise-resume/internal/core/utils"
)

type ResumeHandler struct {
	resumeService ports.ResumeService
	userPort      ports.UserPort
	tokenPort     ports.TokenService
}

func NewResumeHandler(
	resumeService ports.ResumeService,
	userPort ports.UserPort,
	tokenPort ports.TokenService,

) *ResumeHandler {
	return &ResumeHandler{
		resumeService: resumeService,
		userPort:      userPort,
		tokenPort:     tokenPort,
	}
}

func (h ResumeHandler) RegisterResumeRoutes(router fiber.Router) {
	authRouter := router.Group("/resume")
	authRouter.Post(
		"/create",
		middleware.ValidationMiddleware(&dto.CreateResumeDto{}),
		middleware.AuthMiddleware(h.tokenPort, h.userPort),
		h.HandleCreateResume,
	)

	authRouter.Get(
		"/list",
		middleware.AuthMiddleware(h.tokenPort, h.userPort),
		h.HandleFindResumes,
	)
}

// Handles the process of creating a new resume
func (h *ResumeHandler) HandleCreateResume(c *fiber.Ctx) error {
	var body dto.CreateResumeDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userId := c.Locals("user_id")

	res, err := h.resumeService.CreateResume(
		context.WithValue(context.Background(), utils.USER_ID_KEY, userId),
		body,
	)
	if err != nil {
		data := utils.FormatApiResponse(
			"Resume creation failed",
			fiber.Map{"error": err.Error()},
		)
		return c.Status(fiber.StatusForbidden).JSON(data)
	}

	data := utils.FormatApiResponse(
		"Resume was created successfully",
		res,
	)
	return c.Status(fiber.StatusCreated).JSON(data)
}

// Handles the process of listing resumes
func (h *ResumeHandler) HandleFindResumes(c *fiber.Ctx) error {
	userId := c.Locals("user_id")
	res, err := h.resumeService.FindResumes(
		context.WithValue(context.Background(), utils.USER_ID_KEY, userId),
		dto.ResumeFilterDto{
			UserId: userId.(string),
		},
	)
	if err != nil {
		data := utils.FormatApiResponse(
			"Unable to find resumes",
			fiber.Map{"error": err.Error()},
		)
		return c.Status(fiber.StatusForbidden).JSON(data)
	}

	data := utils.FormatApiResponse(
		"Resumes obtained successfully",
		res,
	)
	return c.Status(fiber.StatusOK).JSON(data)
}
