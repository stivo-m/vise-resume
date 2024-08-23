package repository

import (
	"context"

	"github.com/stivo-m/vise-resume/internal/adapters/database"
	"github.com/stivo-m/vise-resume/internal/core/domain"
	"github.com/stivo-m/vise-resume/internal/core/dto"
	"gorm.io/gorm"
)

type ResumeRepository struct {
	db *database.DB
}

func NewResumeRepository(db *database.DB) *ResumeRepository {
	return &ResumeRepository{db: db}
}

func (repo ResumeRepository) CreateResume(ctx context.Context, resume dto.ResumeDto) (*domain.Resume, error) {
	payload := domain.Resume{
		UserId:  resume.UserId,
		Summary: resume.Summary,
		Skills:  resume.Skills,
	}
	result := repo.db.Db.Create(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func (repo ResumeRepository) FindResumeById(ctx context.Context, id string) (*domain.Resume, error) {
	var resume domain.Resume
	result := repo.db.Db.Where("id = ?", id).First(&resume)
	if result.Error != nil {
		return nil, result.Error
	}

	return &resume, nil
}

func (repo ResumeRepository) FindResumeList(ctx context.Context, filter dto.ResumeFilterDto) ([]domain.Resume, error) {
	var resumes []domain.Resume
	result := repo.db.Db.Where("user_id = ?", filter.UserId).Limit(10).Find(&resumes)
	if result.Error != nil {
		return nil, result.Error
	}

	return resumes, nil

}
func (repo ResumeRepository) UpdateResume(ctx context.Context, id string, updates map[string]interface{}) error {
	result := repo.db.Db.Model(&domain.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

}
func (repo ResumeRepository) DeleteResume(ctx context.Context, id string) error {
	result := repo.db.Db.Delete(&domain.Resume{Base: domain.Base{ID: id}})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
func (repo ResumeRepository) AddWorkExperiences(ctx context.Context, id string, experiences []dto.WorkExperienceDto) error {

	var records []domain.WorkExperience
	for _, record := range experiences {
		experience := domain.WorkExperience{
			ResumeId:    id,
			CompanyName: record.CompanyName,
			Role:        record.Role,
			StartDate:   record.StartDate,
			EndDate:     record.EndDate,
		}
		records = append(records, experience)
	}

	result := repo.db.Db.Create(&records)

	if result.Error != nil {
		return result.Error
	}

	return nil

}
func (repo ResumeRepository) UpdateWorkExperiences(ctx context.Context, id string, updates map[string]interface{}) error {
	result := repo.db.Db.Model(&domain.WorkExperience{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repo ResumeRepository) AddEducation(ctx context.Context, id string, education []dto.EducationDto) error {

	var records []domain.Education
	for _, record := range education {
		experience := domain.Education{
			ResumeId:   id,
			SchoolName: record.SchoolName,
			Course:     record.Course,
			StartDate:  record.StartDate,
			EndDate:    record.EndDate,
		}
		records = append(records, experience)
	}

	result := repo.db.Db.Create(&records)

	if result.Error != nil {
		return result.Error
	}

	return nil

}
func (repo ResumeRepository) UpdateEducation(ctx context.Context, id string, updates map[string]interface{}) error {
	result := repo.db.Db.Model(&domain.Education{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repo ResumeRepository) DeleteWorkExperience(ctx context.Context, experienceId string) error {
	result := repo.db.Db.Delete(&domain.WorkExperience{Base: domain.Base{ID: experienceId}})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repo ResumeRepository) DeleteEducation(ctx context.Context, educationId string) error {
	result := repo.db.Db.Delete(&domain.Education{Base: domain.Base{ID: educationId}})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
