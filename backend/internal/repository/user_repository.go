package repository

import (
    "valibibe/internal/models"
    "gorm.io/gorm"
)

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
    CreateUser(user *models.User) error
    GetUserByEmail(email string) (*models.User, error)
    GetUserByID(id string) (*models.User, error)
}

// userRepository - конкретная реализация UserRepository
type userRepository struct {
    db *gorm.DB
}

// NewUserRepository создает новый экземпляр репозитория пользователей
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

// CreateUser добавляет нового пользователя в базу данных
func (r *userRepository) CreateUser(user *models.User) error {
    return r.db.Create(user).Error
}

// GetUserByEmail ищет пользователя по email
func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
    var user models.User
    if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

// GetUserByID ищет пользователя по ID
func (r *userRepository) GetUserByID(id string) (*models.User, error) {
    var user models.User
    if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
    		return nil, err
    	}
    return &user, nil
}











