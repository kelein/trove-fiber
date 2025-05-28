package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/kelein/trove-fiber/internal/model"
)

// UserRepository abstracts the user-related operations
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepository struct{ *Repository }

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(r *Repository) UserRepository {
	return &userRepository{Repository: r}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.DB(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.DB(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	cond := map[string]string{"user_id": userID}
	if err := r.DB(ctx).Where(cond).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	cond := map[string]string{"email": email}
	if err := r.DB(ctx).Where(cond).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
