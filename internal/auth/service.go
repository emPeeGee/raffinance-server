package auth

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/emPeeGee/raffinance/pkg/crypt"
	"github.com/emPeeGee/raffinance/pkg/log"
	"github.com/golang-jwt/jwt/v4"
)

const (
	signingKey = "bv646gf930ds^#fg&)Fd_)))*("
	tokenTTL   = time.Hour * 12
)

type Service interface {
	createUser(input createUserDTO) error
	generateToken(credentials credentialsDTO) (string, error)
	getUserById(id uint) (UserResponse, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId uint `json:"userId"`
}

func NewAuthService(repository Repository, logger log.Logger) *service {
	return &service{repo: repository, logger: logger}
}

func (s *service) createUser(user createUserDTO) error {
	// TODO crypt
	hashedPassword, err := crypt.HashPassphrase(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	err = s.repo.createUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) generateToken(credentials credentialsDTO) (string, error) {
	hashedPassword, err := s.repo.getHashedPasswordByUsername(credentials.Username)
	if err != nil {
		return "", err
	}

	ok := crypt.CheckPasswordHashes(credentials.Password, hashedPassword.Password)
	if !ok {
		return "", errors.New("password does not match")
	}

	user, err := s.repo.getUserByUsername(credentials.Username)
	if err != nil {
		return "", err
	}

	if err := s.repo.updateLatestLogins(credentials.Username); err != nil {
		return "", err
	}

	str, _ := json.MarshalIndent(user, "", "\t")
	s.logger.Debug(string(str))

	// TODO: constants to be moved in config
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *service) getUserById(id uint) (UserResponse, error) {
	return s.repo.getUserById(id)
}
