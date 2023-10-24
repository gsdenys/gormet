package gormet

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Repository[T any] struct {
	DB     *gorm.DB
	Config *Config
}

func CreateRepository[T any](db *gorm.DB, conf *Config) *Repository[T] {
	if conf == nil {
		conf = DefaultConfig()
	}

	return &Repository[T]{
		DB:     db,
		Config: conf,
	}
}

func New[T any](db *gorm.DB) *Repository[T] {
	return (*Repository[T])(CreateRepository[T](db, nil))
}

func (r *Repository[T]) Create(entity *T) error {
	if entity == nil {
		return errors.New("the entity cannot be nil")
	}

	if err := isValid[T](entity, r.Config.Validate); err != nil {
		return err
	}

	if result := r.DB.Create(entity); result.Error != nil {
		return result.Error
	}

	return nil
}

func isValid[T any](entity *T, validable bool) error {
	if validable {
		validate := validator.New()
		return validate.Struct(entity)
	}

	return nil
}
