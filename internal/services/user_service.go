package services

import (
	"errors"
	"fmt"
	"math/rand"
	"pedika-layered-architecture/internal/models"
	"pedika-layered-architecture/internal/repositories"
	"pedika-layered-architecture/internal/utils"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *models.User, confirmPassword string) error
	Login(email, password string) (*LoginResponse, error)
}

type LoginResponse struct {
	Token       string `json:"token"`
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(user *models.User, confirmPassword string) error {
	if user.Password != confirmPassword {
		return errors.New("password and confirm password tidak sama")
	}

	if s.repo.IsEmailExists(user.Email) {
		return errors.New("email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	firstWord := strings.Split(user.FullName, " ")[0]
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(100000)
	user.Username = strings.ToLower(firstWord) + fmt.Sprintf("%05d", randomNumber)

	user.Role = "masyarakat"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.repo.Create(user)
}

func (s *userService) Login(email, password string) (*LoginResponse, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("email tidak ditemukan")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("password salah")
	}

	token, err := utils.GenerateToken(user.ID, user.Role, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:       token,
		ID:          user.ID,
		Username:    user.Username,
		FullName:    user.FullName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
	}, nil
}
