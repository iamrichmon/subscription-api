package repository

import (
	"errors"

	"github.com/iamrichmon/subscription-api/internal/model"
	"github.com/iamrichmon/subscription-api/internal/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) Create(admin *model.Admin) error {
	err := r.db.Create(admin).Error
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
		return utils.ErrEmailTaken
	}

	return err
}

func (r *AdminRepository) FindByEmail(email string) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepository) FindAll() ([]model.User, error) {
	users := make([]model.User, 0)

	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
