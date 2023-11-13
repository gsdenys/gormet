package gormet

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRepository_GetById(t *testing.T) {

	type testGet struct {
		gorm.Model
		Email string `json:"email" gorm:"unique,not null;default:null"`
		Name  string `json:"name" gorm:"unique;not null;default:null"`
	}

	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, db)

	db.AutoMigrate(&testGet{})

	repo, err := New[testGet](db)
	assert.Nil(t, err)

	createRegister := func() *testGet {
		var field = uuid.NewString()
		entity := &testGet{
			Name:  fmt.Sprintf("%v", field),
			Email: fmt.Sprintf("%v@mail.com", field),
		}

		repo.Create(entity)

		return entity
	}

	t.Run("Get Successfull", func(t *testing.T) {
		entity := createRegister()
		id := entity.ID

		got, err := repo.GetById(id)

		assert.NotNil(t, got)
		assert.Nil(t, err)
		assert.Equal(t, id, got.ID)
	})

	t.Run("Nil entity", func(t *testing.T) {
		got, err := repo.GetById(nil)

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, "the id should not be nil", err.Error())
	})

	t.Run("Entity with sql injection", func(t *testing.T) {
		entity := createRegister()

		id := fmt.Sprintf(`%d"; update test_deletes set deleted at '1970-01-01 00:00:00 --`, entity.ID)

		got, err := repo.GetById(id)

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, "record not found", err.Error())
	})

	t.Run("Connection Closed", func(t *testing.T) {
		entity := createRegister()

		sqlDB, _ := db.DB()
		err := sqlDB.Close()
		assert.Nil(t, err)

		got, err := repo.GetById(entity.ID)

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "sql: database is closed")
	})
}
