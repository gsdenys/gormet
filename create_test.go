package gormet

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRepository_Create(t *testing.T) {

	type testCreate struct {
		gorm.Model
		Email string `json:"email" gorm:"unique,not null;default:null"`
		Name  string `json:"name" gorm:"unique;not null;default:null"`
	}

	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, db)

	db.AutoMigrate(&testCreate{})

	repo, err := New[testCreate](db)
	assert.Nil(t, err)

	t.Run("Creation Successfull", func(t *testing.T) {
		field := uuid.NewString()
		entity := &testCreate{
			Name:  fmt.Sprintf("create-%v", field),
			Email: fmt.Sprintf("%v@mail.com", field),
		}

		err := repo.Create(entity)
		assert.Nil(t, err)

		assert.Greater(t, entity.ID, uint(0))
	})

	t.Run("Nil entity", func(t *testing.T) {
		err := repo.Create(nil)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "the entity should not be nil")
	})

	t.Run("Entity with sql injection", func(t *testing.T) {
		field := uuid.NewString()

		//The intire string should be placed in tha field name
		injection := fmt.Sprintf(`sqli-%s"); update test_creates set deleted at '1970-01-01 00:00:00' -- `, field)

		entity := &testCreate{
			Email: fmt.Sprintf("%v@mail.com", field),
			Name:  injection,
		}

		err = repo.Create(entity)
		assert.Nil(t, err)

		got, err := repo.GetById(entity.ID)
		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, injection, got.Name)
	})

	t.Run("Connection Closed", func(t *testing.T) {
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
}
