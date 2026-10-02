package service

import (
	"errors"

	"strings"

	"github.com/iamrichmon/subscription-api/internal/auth"
	"github.com/iamrichmon/subscription-api/internal/model"
	"github.com/iamrichmon/subscription-api/internal/repository"
	"github.com/iamrichmon/subscription-api/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminService struct {
	repo      *repository.AdminRepository
	jwtSecret string
}

func NewAdminService(repo *repository.AdminRepository, jwtSecret string) *AdminService {
	return &AdminService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *AdminService) RegisterAdmin(name, email, password string) (*model.Admin, error) {
	// Normalize name and email
	name = utils.NormalizeName(name)

	if name == "" || password == "" {
		return nil, errors.New("Name and password are required.")
	}

	// email validation
	addr, err := utils.NormalizeEmail(email)

	if err != nil {
		return nil, err
	}

	email = addr

	// password hashed
	hashed, err := utils.HashPlainPass(password)
	if err != nil {
		return nil, err
	}

	admin := &model.Admin{
		Name:     name,
		Email:    email,
		Password: hashed,
		Role:     model.RoleAdmin,
	}

	if err := s.repo.Create(admin); err != nil {
		return nil, err
	}

	return admin, nil
}

func (s *AdminService) Login(email, password string) (string, error) {

	if strings.TrimSpace(password) == "" {
		return "", utils.ErrInvalidCredentials
	}

	addr, err := utils.NormalizeEmail(email)

	if err != nil {
		return "", utils.ErrInvalidCredentials
	}

	email = addr

	existingAdmin, err := s.repo.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err // DB/infra err, e.g. connection error
	}

	if existingAdmin == nil {
		return "", utils.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingAdmin.Password), []byte(password)); err != nil {
		return "", utils.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(existingAdmin.ID, existingAdmin.Email, existingAdmin.Role, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AdminService) GetAllUsers() ([]model.User, error) {
	users, err := s.repo.FindAll()

	if err != nil {
		return nil, err
	}

	return users, nil
}
