package service

import (
	"errors"
	"time"

	"mini_url/internal/models"
	"mini_url/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(username, email, password string) (*models.User, error) {
	existing, _ := s.userRepo.FindByUsername(username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	existingEmail, _ := s.userRepo.FindByEmail(email)
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	user := &models.User{
		Username:  username,
		Email:     email,
		Password:  password,
		Role:      "user",
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if user.Password != password {
		return "", errors.New("invalid username or password")
	}

	user.LastLogin = time.Now()
	s.userRepo.Update(user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}