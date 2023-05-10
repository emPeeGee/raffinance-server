package auth

import (
	"time"

	"github.com/emPeeGee/raffinance/internal/entity"
	"github.com/emPeeGee/raffinance/pkg/log"
	"gorm.io/gorm"
)

type Repository interface {
	createUser(user createUserDTO) error
	updateLatestLogins(username string) error
	getUserById(id uint) (UserResponse, error)
	getUserByUsername(username string) (entity.User, error)
	getHashedPasswordByUsername(username string) (userHashedPassword, error)
}

type repository struct {
	db     *gorm.DB
	logger log.Logger
}

func NewAuthRepository(db *gorm.DB, logger log.Logger) *repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) createUser(user createUserDTO) error {
	newUser := entity.User{
		Username:     user.Username,
		Password:     user.Password,
		Name:         user.Name,
		Email:        user.Email,
		Phone:        user.Phone,
		LatestLogins: []string{},
	}

	if err := r.db.Create(&newUser).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) getUserByUsername(username string) (entity.User, error) {
	var user entity.User

	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *repository) updateLatestLogins(username string) error {
	var user entity.User

	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}

	// NOTE: Now().String() will include monotonic clock
	user.LatestLogins = append(user.LatestLogins, time.Now().Format(time.RFC3339))
	if err := r.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) getHashedPasswordByUsername(username string) (userHashedPassword, error) {
	var userPassword userHashedPassword

	if err := r.db.Model(&entity.User{}).Where("username = ?", username).First(&userPassword).Error; err != nil {
		return userHashedPassword{}, err
	}

	return userPassword, nil
}

func (r *repository) getUserById(id uint) (UserResponse, error) {
	var user UserResponse

	if err := r.db.Model(&entity.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return UserResponse{}, err
	}

	r.logger.Debug(user)

	return user, nil
}
