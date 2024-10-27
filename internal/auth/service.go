package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/courage173/go-auth-api/internal/models"

	"github.com/courage173/go-auth-api/internal/users"

	"github.com/courage173/go-auth-api/pkg/log"

	"golang.org/x/crypto/bcrypt"

	"strings"

	"github.com/courage173/go-auth-api/internal/errors"

	"github.com/dgrijalva/jwt-go"
)


type Service interface {
	Register(ctx context.Context, user models.User) (string, error)
	Login(ctx context.Context, email string, password string) (string, error)
}


type service struct {
	signingKey string
	tokenExpiration int
	logger log.Logger
	userRepo   users.Repository
}

type TokenIdentity interface {
	// GetID returns the user ID.
	GetID() int
	// GetName returns the user name.
	GetEmail() string
}

func NewService(signingKey string, tokenExpiration int, logger log.Logger, userRepo users.Repository) Service {
	return service{
        signingKey,
        tokenExpiration,
        logger,
        userRepo,
    }
}

func (s service) Register(ctx context.Context, data models.User) (string, error) {
	// Implementation of user registration logic
    // Here, we assume the user is stored in a database and hashed password is stored securely

    // token, err := generateToken(user.ID, email, s.signingKey, s.tokenExpiration)
    // if err!= nil {
    //     s.logger.Error(ctx, "Error generating token", "error", err)
    //     return "", err
    // }

    // return token, nil

	//lowercase the email address
	email := strings.ToLower(data.Email)

	userExist := s.userRepo.EmailExist(ctx, email)

	if userExist {
        return "", errors.BadRequest("User already registered")
    }

	
	passwordHashed, err := hashedPassword(data.Password)
	if err!= nil {
        s.logger.Error(ctx, "Error hashing password", "error", err)
        return "", err
    }

	// save both hased password and lowercased email
	data.Password =  passwordHashed
	data.Email = email

	 error := s.userRepo.Create(ctx, data)

	 if error!= nil {
         s.logger.Error(ctx, "Error creating user", "error", error)
         return "", error
     }

	 return "Account created successfully", nil
}

func (s service) Login(ctx context.Context, email string, password string) (string, error) {
	email = strings.ToLower(email)
    fmt.Println( s.userRepo)
    user, err := s.userRepo.GetByEmail(ctx, email)
    if err!= nil {
        return "", err
    }

    if user.ID == 0 {
        return "", errors.NotFound("User not found")
    }

    isMatch, err := verifyPassword(password, user.Password)
    if err!= nil {
        s.logger.Error(ctx, "Error verifying password", "error", err)
        return "", err
    }

    if!isMatch {
        return "", errors.BadRequest("Invalid password")
    }

    return s.generateJWT(models.User{ID: user.ID, Email: email})
}

func hashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err!= nil {
        return "", err
    }
    return string(hashedPassword), nil
}

func verifyPassword(password, hashedPassword string) (bool, error){
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err == bcrypt.ErrMismatchedHashAndPassword {
        return false, nil
    } else if err!= nil {
        return false, err
    }
    return true, nil
}

func (s service) generateJWT(identity TokenIdentity) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   identity.GetID(),
		"email": identity.GetEmail(),
		"exp":  time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(s.signingKey))
}