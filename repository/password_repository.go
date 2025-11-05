package repository

import (
    "github.com/google/uuid"
    "github.com/mhandyalf/go-passmanager/models"
    "gorm.io/gorm"
)

type PasswordRepository interface {
    Create(p *models.Password) error
    GetByUserID(userID uuid.UUID) ([]models.Password, error)
    GetByID(id string) (*models.Password, error)
    Update(p *models.Password, updates map[string]interface{}) error
    Delete(p *models.Password) error
}

type passwordRepo struct {
    db *gorm.DB
}

func NewPasswordRepository(db *gorm.DB) PasswordRepository {
    return &passwordRepo{db: db}
}

func (r *passwordRepo) Create(p *models.Password) error {
    return r.db.Create(p).Error
}

func (r *passwordRepo) GetByUserID(userID uuid.UUID) ([]models.Password, error) {
    var list []models.Password
    if err := r.db.Where("user_id = ?", userID).Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}

func (r *passwordRepo) GetByID(id string) (*models.Password, error) {
    var p models.Password
    if err := r.db.First(&p, id).Error; err != nil {
        return nil, err
    }
    return &p, nil
}

func (r *passwordRepo) Update(p *models.Password, updates map[string]interface{}) error {
    return r.db.Model(p).Updates(updates).Error
}

func (r *passwordRepo) Delete(p *models.Password) error {
    return r.db.Delete(p).Error
}
