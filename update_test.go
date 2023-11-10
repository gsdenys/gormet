package gormet

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRepository_Update(t *testing.T) {

	type testUpdate struct {
		gorm.Model
		Email string `json:"email" gorm:"unique,not null;default:null"`
		Name  string `json:"name" gorm:"unique;not null;default:null"`
	}

	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, db)

	db.AutoMigrate(&testUpdate{})

	repo, err := New[testUpdate](db)
	assert.Nil(t, err)

	createRegister := func() *testUpdate {
		var field = uuid.NewString()
		entity := &testUpdate{
			Name:  fmt.Sprintf("%v", field),
			Email: fmt.Sprintf("%v@mail.com", field),
		}

		repo.Create(entity)

		return entity
	}

	t.Run("Updated Successfull", func(t *testing.T) {
		entity := createRegister()
		newValue := fmt.Sprintf("updated-%s", uuid.NewString())

		entity.Name = newValue
		repo.Update(entity)

		assert.Greater(t, entity.ID, uint(0))
	})

	t.Run("Nil entity", func(t *testing.T) {
		err := repo.Update(nil)

		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "the entity should not be nil")
	})

	t.Run("Entity with sql injection", func(t *testing.T) {
		entity := createRegister()

		field := uuid.NewString()
		injection := fmt.Sprintf(`sqli-%s\"); update test_creates set deleted at '1970-01-01 00:00:00' -- `, field)
		entity.Name = injection

		err = repo.Update(entity)
		assert.Nil(t, err)

		got, err := repo.GetById(entity.ID)
		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, injection, got.Name)
	})

	t.Run("Connection Closed", func(t *testing.T) {
		entity := createRegister()

		sqlDB, _ := db.DB()
		err := sqlDB.Close()
		assert.Nil(t, err)

		entity.Name = uuid.NewString()

		err = repo.Update(entity)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "sql: database is closed")
	})
}
