package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/yourusername/user-api/internal/models"
	"github.com/yourusername/user-api/internal/repository"
)

const dobLayout = "2006-01-02"

// ErrUserNotFound is returned when a user does not exist in the database.
var ErrUserNotFound = errors.New("user not found")

// UserService defines all business-logic operations for the user domain.
type UserService interface {
	Create(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error)
	GetByID(ctx context.Context, id uint32) (*models.UserWithAgeResponse, error)
	Update(ctx context.Context, id uint32, req *models.UpdateUserRequest) (*models.UserResponse, error)
	Delete(ctx context.Context, id uint32) error
	List(ctx context.Context, page, limit int32) (*models.ListUsersResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService returns a UserService backed by the given repository.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// -----------------------------------------------------------------------
// CalculateAge computes the age from a date of birth to today.
// This is the single source of truth for age calculation in the service layer.
// -----------------------------------------------------------------------
func CalculateAge(dob time.Time) int {
	today := time.Now().UTC()
	age := today.Year() - dob.Year()
	if today.Month() < dob.Month() ||
		(today.Month() == dob.Month() && today.Day() < dob.Day()) {
		age--
	}
	return age
}

// -----------------------------------------------------------------------
// Create
// -----------------------------------------------------------------------

func (s *userService) Create(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	dob, err := time.Parse(dobLayout, req.DOB)
	if err != nil {
		return nil, err
	}

	id, err := s.repo.Create(ctx, req.Name, dob)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:   uint32(id),
		Name: req.Name,
		DOB:  req.DOB,
	}, nil
}

// -----------------------------------------------------------------------
// GetByID
// -----------------------------------------------------------------------

func (s *userService) GetByID(ctx context.Context, id uint32) (*models.UserWithAgeResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &models.UserWithAgeResponse{
		ID:   u.ID,
		Name: u.Name,
		DOB:  u.Dob.Format(dobLayout),
		Age:  CalculateAge(u.Dob),
	}, nil
}

// -----------------------------------------------------------------------
// Update
// -----------------------------------------------------------------------

func (s *userService) Update(ctx context.Context, id uint32, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	// Verify the user exists before updating.
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	dob, err := time.Parse(dobLayout, req.DOB)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, id, req.Name, dob); err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:   id,
		Name: req.Name,
		DOB:  req.DOB,
	}, nil
}

// -----------------------------------------------------------------------
// Delete
// -----------------------------------------------------------------------

func (s *userService) Delete(ctx context.Context, id uint32) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}

// -----------------------------------------------------------------------
// List
// -----------------------------------------------------------------------

func (s *userService) List(ctx context.Context, page, limit int32) (*models.ListUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	data := make([]models.UserWithAgeResponse, 0, len(users))
	for _, u := range users {
		data = append(data, models.UserWithAgeResponse{
			ID:   u.ID,
			Name: u.Name,
			DOB:  u.Dob.Format(dobLayout),
			Age:  CalculateAge(u.Dob),
		})
	}

	return &models.ListUsersResponse{
		Data:  data,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}
