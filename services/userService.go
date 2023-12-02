package services

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xbt573/project-example/models"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserServiceAlreadyExists = errors.New("user already exists")
	ErrUserServiceUnauthorized  = errors.New("unauthorized")
)

type UserService interface {
	Register(string, string) (models.Tokens, error)
	Login(string, string) (models.Tokens, error)
	Refresh(string) (models.Tokens, error)
}

type concreteUserService struct {
	secret   string
	database *gorm.DB
}

func NewUserService(database *gorm.DB, secret string) UserService {
	return &concreteUserService{secret, database}
}

func (s *concreteUserService) Register(login, password string) (models.Tokens, error) {
	var user models.User
	result := s.database.Where("login = ?", login).First(&user)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Tokens{}, result.Error
		}
	}

	if result.RowsAffected > 0 {
		return models.Tokens{}, ErrUserServiceAlreadyExists
	}

	sha := sha256.New()
	sha.Write([]byte(password))
	hash := fmt.Sprintf("%x", sha.Sum(nil))

	user = models.User{
		Login:        login,
		PasswordHash: hash,
	}
	result = s.database.Create(&user)
	if result.Error != nil {
		return models.Tokens{}, result.Error
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.ID,
		"subject": "refresh_token",
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.ID,
		"subject": "access_token",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	refreshSigned, err := refreshToken.SignedString([]byte(s.secret))
	if err != nil {
		return models.Tokens{}, err
	}
	accessSigned, err := accessToken.SignedString([]byte(s.secret))
	if err != nil {
		return models.Tokens{}, err
	}

	return models.Tokens{
		RefreshToken: refreshSigned,
		AccessToken:  accessSigned,
	}, nil
}

func (s *concreteUserService) Login(login, password string) (models.Tokens, error) {
	var user models.User

	result := s.database.Where("login = ?", login).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Tokens{}, ErrUserServiceUnauthorized
		}

		return models.Tokens{}, result.Error
	}

	sha := sha256.New()
	sha.Write([]byte(password))
	hash := fmt.Sprintf("%x", sha.Sum(nil))

	if user.Login != login || user.PasswordHash != hash {
		return models.Tokens{}, ErrUserServiceUnauthorized
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.ID,
		"subject": "refresh_token",
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.ID,
		"subject": "access_token",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	refreshSigned, err := refreshToken.SignedString([]byte(s.secret))
	if err != nil {
		return models.Tokens{}, err
	}
	accessSigned, err := accessToken.SignedString([]byte(s.secret))
	if err != nil {
		return models.Tokens{}, err
	}

	return models.Tokens{
		RefreshToken: refreshSigned,
		AccessToken:  accessSigned,
	}, nil
}

func (s *concreteUserService) Refresh(refreshToken string) (models.Tokens, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return models.Tokens{}, errors.New("something got wrong i guess 🐳")
		}

		return []byte(s.secret), nil
	})
	if err != nil {
		return models.Tokens{}, err
	}

	claims := token.Claims.(jwt.MapClaims)
	userid := uint(claims["id"].(float64))

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      userid,
		"subject": "refresh_token",
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      userid,
		"subject": "access_token",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	newRefreshSigned, err := newRefreshToken.SignedString([]byte(s.secret))
	if err != nil {
		return models.Tokens{}, err
	}
	newAccessSigned, err := newAccessToken.SignedString([]byte(s.secret))
	if err != nil {
		return models.Tokens{}, err
	}

	return models.Tokens{
		RefreshToken: newRefreshSigned,
		AccessToken:  newAccessSigned,
	}, nil
}
