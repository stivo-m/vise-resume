package services

import (
	"context"

	"github.com/stivo-m/vise-resume/internal/core/dto"
	"github.com/stivo-m/vise-resume/internal/core/ports"
	"github.com/stivo-m/vise-resume/internal/core/utils"
)

type ResumeService struct {
	resumePort ports.ResumePort
}

func NewResumeService(
	resumePort ports.ResumePort,

) *ResumeService {
	return &ResumeService{
		resumePort: resumePort,
	}
}

func (s ResumeService) CreateResume(ctx context.Context, payload dto.CreateResumeDto) (*dto.ResumeDto, error) {
	userId, err := utils.ExtractUuidFromContext(ctx, utils.USER_ID_KEY)
	if err != nil {
		return nil, err
	}

	resume, err := s.resumePort.CreateResume(ctx, dto.ResumeDto{
		UserId:  userId.String(),
		Summary: payload.Summary,
		Skills:  payload.Skills,
	})

	if err != nil {
		return nil, err
	}

	err = s.resumePort.AddEducation(ctx, resume.ID, payload.Education)
	if err != nil {
		return nil, err
	}

	err = s.resumePort.AddWorkExperiences(ctx, resume.ID, payload.Experiences)
	if err != nil {
		return nil, err
	}

	return &dto.ResumeDto{
		ID:      resume.ID,
		UserId:  userId.String(),
		Summary: resume.Summary,
		Skills:  resume.Skills,
	}, nil
}

func (s ResumeService) FindResumes(ctx context.Context, payload dto.ResumeFilterDto) ([]dto.ResumeDto, error) {
	resumes, err := s.resumePort.FindResumeList(ctx, payload)
	if err != nil {
		return nil, err
	}

	var result []dto.ResumeDto
	for _, resume := range resumes {
		result = append(result, dto.ResumeDto{
			ID:      resume.ID,
			UserId:  payload.UserId,
			Summary: resume.Summary,
			Skills:  resume.Skills,
		})
	}

	return result, nil
}
