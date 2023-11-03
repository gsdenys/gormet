package gormet

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testCreate struct {
	gorm.Model
	Name  string `json:"name" gorm:"unique;not null;default:null" validate:"required,min=3"`
	Email string `json:"email" gorm:"unique,not null;default:null" validate:"required,email"`
}

func TestRepository_Create(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, db)

	db.AutoMigrate(&testCreate{})

	repo := New[testCreate](db)

	t.Run("Create entity successfull", func(t *testing.T) {
		field := uuid.NewString()
		entity := &testCreate{
			Name:  fmt.Sprintf("create-%v", field),
			Email: fmt.Sprintf("%v@mail.com", field),
		}

		err := repo.Create(entity)
		assert.Nil(t, err)

		assert.Greater(t, entity.ID, uint(0))
	})

	t.Run("Create nil entity error", func(t *testing.T) {
		err := repo.Create(nil)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "The entity cannot be nil")
	})

	t.Run("Create validation error", func(t *testing.T) {
		field := uuid.NewString()
		entity := &testCreate{
			Name:  field,
			Email: field,
		}

		err := repo.Create(entity)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "Key: 'testCreate.Email' Error:Field validation for 'Email' failed on the 'email' tag")
	})

	t.Run("Create entity with sql injection", func(t *testing.T) {
		field := uuid.NewString()
		entity := &testCreate{
			Name:  "sqlInjection" + field + "; update test_creates set deleted at '1970-01-01 00:00:00'",
			Email: fmt.Sprintf("%v@mail.com", field),
		}

		err = repo.Create(entity)
		assert.Nil(t, err)

		got, err := repo.GetById(entity.ID)
		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "sqlInjection"+field+"; update test_creates set deleted at '1970-01-01 00:00:00'", got.Name)
	})

	t.Run("Create close connection error", func(t *testing.T) {
		sqlDB, _ := db.DB()
		err := sqlDB.Close()
		assert.Nil(t, err)

		field := uuid.NewString()
		entity := &testCreate{
			Name:  field,
			Email: fmt.Sprintf("%v@mail.com", field),
		}

		err = repo.Create(entity)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "sql: database is closed")
	})

	t.Run("Create entity successfull without validation", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
		assert.Nil(t, err)
		assert.NotNil(t, db)

		repo := CreateRepository[testCreate](db, &Config{
			Validate: false,
		})

		field := uuid.NewString()
		entity := &testCreate{
			Name:  fmt.Sprintf("noValidation-%v", field),
			Email: field,
		}

		err = repo.Create(entity)
		assert.Nil(t, err)

		assert.Greater(t, entity.ID, uint(0))
	})

}
