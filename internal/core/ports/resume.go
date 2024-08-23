package ports

import (
	"context"

	"github.com/stivo-m/vise-resume/internal/core/domain"
	"github.com/stivo-m/vise-resume/internal/core/dto"
)

type ResumePort interface {
	CreateResume(ctx context.Context, resume dto.ResumeDto) (*domain.Resume, error)
	FindResumeById(ctx context.Context, id string) (*domain.Resume, error)
	FindResumeList(ctx context.Context, filter dto.ResumeFilterDto) ([]domain.Resume, error)
	UpdateResume(ctx context.Context, id string, updates map[string]interface{}) error
	DeleteResume(ctx context.Context, id string) error
	AddWorkExperiences(ctx context.Context, id string, experiences []dto.WorkExperienceDto) error
	UpdateWorkExperiences(ctx context.Context, id string, updates map[string]interface{}) error
	AddEducation(ctx context.Context, id string, education []dto.EducationDto) error
	UpdateEducation(ctx context.Context, id string, updates map[string]interface{}) error
	DeleteWorkExperience(ctx context.Context, experienceId string) error
	DeleteEducation(ctx context.Context, educationId string) error
}

type ResumeService interface {
	CreateResume(ctx context.Context, payload dto.CreateResumeDto) (*dto.ResumeDto, error)
	FindResumes(ctx context.Context, payload dto.ResumeFilterDto) ([]dto.ResumeDto, error)
}
