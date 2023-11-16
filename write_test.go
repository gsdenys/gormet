package gormet

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type testWrite struct {
	gorm.Model
	Email string `json:"email" gorm:"unique,not null;default:null"`
	Name  string `json:"name" gorm:"unique;not null;default:null"`
}

func TestRepository_Create(t *testing.T) {

	db := getGormConnection(t, &testWrite{})

	repo, err := New[testWrite](db)
	assert.Nil(t, err)

	t.Run("Create entity successfully", func(t *testing.T) {
		field := uuid.NewString()
		entity := &testWrite{
			Name:  fmt.Sprintf("create-%v", field),
			Email: fmt.Sprintf("%v@mail.com", field),
		}

		err := repo.Create(entity)
		assert.Nil(t, err)

		assert.Greater(t, entity.ID, uint(0))
	})

	t.Run("Error send a nil entity", func(t *testing.T) {
		err := repo.Create(nil)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "the entity should not be nil")
	})

	t.Run("Entity with sql injection", func(t *testing.T) {
		field := uuid.NewString()

		//The intire string should be placed in tha field name
		injection := fmt.Sprintf("sql-%s", field) + `"; update test_creates set deleted at '1970-01-01 00:00:00' -- `

		entity := &testWrite{
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

	t.Run("Error connection closed", func(t *testing.T) {
		sqlDB, _ := db.DB()
		err := sqlDB.Close()
		assert.Nil(t, err)

		field := uuid.NewString()
		entity := &testWrite{
			Name:  field,
			Email: fmt.Sprintf("%v@mail.com", field),
		}

		err = repo.Create(entity)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "sql: database is closed")
	})
}

func TestRepository_Update(t *testing.T) {

	db := getGormConnection(t, &testWrite{})

	repo, err := New[testWrite](db)
	assert.Nil(t, err)

	t.Run("Update entity successfully", func(t *testing.T) {
		entity := &testWrite{
			Name:  uuid.NewString(),
			Email: fmt.Sprintf("%s@mail.com", uuid.NewString()),
		}
		err := repo.Create(entity)
		assert.Nil(t, err)

		newValue := fmt.Sprintf("updated-%s", uuid.NewString())
		entity.Name = newValue

		err = repo.Update(entity)
		assert.Nil(t, err)

		got, err := repo.GetById(entity.ID)
		assert.Nil(t, err)

		assert.Equal(t, newValue, got.Name)
	})

	t.Run("Error send Nil entity", func(t *testing.T) {
		err := repo.Update(nil)

		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "the entity should not be nil")
	})

	t.Run("Entity with sql injection", func(t *testing.T) {
		entity := &testWrite{
			Name:  uuid.NewString(),
			Email: fmt.Sprintf("%s@mail.com", uuid.NewString()),
		}
		err := repo.Create(entity)
		assert.Nil(t, err)

		field := uuid.NewString()
		injection := fmt.Sprintf("sql-%s", field) + `update test_updates set deleted at '1970-01-01 00:00:00' -- `
		entity.Name = injection

		err = repo.Update(entity)
		assert.Nil(t, err)

		got, err := repo.GetById(entity.ID)
		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, injection, got.Name)
	})

	t.Run("Error connection closed", func(t *testing.T) {
		entity := &testWrite{
			Name:  uuid.NewString(),
			Email: fmt.Sprintf("%s@mail.com", uuid.NewString()),
		}
		err := repo.Create(entity)
		assert.Nil(t, err)

		sqlDB, _ := db.DB()
		err = sqlDB.Close()
		assert.Nil(t, err)

		entity.Name = uuid.NewString()

		err = repo.Update(entity)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "sql: database is closed")
	})
}
