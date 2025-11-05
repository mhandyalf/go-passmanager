package repository

import (
	"github.com/mhandyalf/go-passmanager/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(u *models.User) error
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	Update(u *models.User) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *userRepo) GetByUsername(username string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("user_name = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByID(id string) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Update(u *models.User) error {
	return r.db.Save(u).Error
}
