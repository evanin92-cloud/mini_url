package repository

import "mini_url/internal/models"

type LinkRepositoryInterface interface {
	Create(link *models.Link) error
	FindByID(id int) (*models.Link, error)
	FindByShortID(shortID string) (*models.Link, error)
	FindByUserID(userID int) ([]*models.Link, error)
	Update(link *models.Link) error
	Delete(id int) error
}

type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByID(id int) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

type StatsRepositoryInterface interface {
	Create(stats *models.Stats) error
	FindByShortID(shortID string) ([]*models.Stats, error)
}