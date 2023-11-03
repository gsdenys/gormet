package gormet

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testUpdate struct {
	gorm.Model
	Name  string `json:"name" gorm:"unique;not null;default:null" validate:"required,min=3"`
	Email string `json:"email" gorm:"unique,not null;default:null" validate:"required,email"`
}

func createUpdateEntity(repo *Repository[testUpdate]) *testUpdate {
	field := uuid.NewString()
	entity := &testUpdate{
		Name:  field,
		Email: fmt.Sprintf("%v@mail.com", field),
	}

	if err := repo.Create(entity); err != nil {
		return nil
	}

	return entity
}

func TestRepository_Update(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, db)

	db.AutoMigrate(&testUpdate{})

	repo := New[testUpdate](db)

	t.Run("Update entity successfull", func(t *testing.T) {
		entity := createUpdateEntity(repo)
		assert.NotNil(t, entity)

		// id := elements[0].ID
		name := entity.Name
		email := entity.Email

		entity.Name = uuid.NewString()
		entity.Email = fmt.Sprintf("%v@mail.com", entity.Name)

		err := repo.Update(entity)
		assert.Nil(t, err)

		updated, err := repo.GetById(entity.ID)
		assert.Nil(t, err)
		assert.NotNil(t, updated)

		assert.NotEqual(t, updated.Name, name)
		assert.NotEqual(t, updated.Email, email)

		assert.Equal(t, int(updated.ID), int(entity.ID))
		assert.Equal(t, updated.Name, entity.Name)
		assert.Equal(t, updated.Email, entity.Email)
	})

	t.Run("Update error nil entity", func(t *testing.T) {
		err := repo.Update(nil)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "the entity cannot be nil")
	})

	t.Run("Update error invalid entity", func(t *testing.T) {
		entity := createUpdateEntity(repo)
		assert.NotNil(t, entity)

		id := entity.ID

		entity.Name = uuid.NewString()
		entity.Email = uuid.NewString() //fail in validation

		err := repo.Update(entity)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "Key: 'testUpdate.Email' Error:Field validation for 'Email' failed on the 'email' tag")

		updated, err := repo.GetById(id)
		assert.Nil(t, err)
		assert.NotNil(t, updated)
	})

	t.Run("Update entity with sql injection", func(t *testing.T) {
		field := uuid.NewString()
		entity := &testUpdate{
			Name:  uuid.NewString(),
			Email: fmt.Sprintf("%v@mail.com", field),
		}

		err = repo.Create(entity)
		assert.Nil(t, err)

		got, err := repo.GetById(entity.ID)
		assert.Nil(t, err)
		assert.NotNil(t, got)

		entity.Name = "sqlInjection-" + field + "\"; update test_creates set deleted at '2023-11-03 15:53:52.606' --"

		err = repo.Update(entity)

		got, err = repo.GetById(entity.ID)

		assert.Nil(t, err)
		assert.Equal(t, entity.Name, got.Name)
	})

	t.Run("Update invalid entity not validated", func(t *testing.T) {
		entity := createUpdateEntity(repo)
		assert.NotNil(t, entity)

		id := entity.ID

		entity.Name = uuid.NewString()
		entity.Email = uuid.NewString() //fail in validation

		db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
		assert.Nil(t, err)
		assert.NotNil(t, db)

		repo2 := CreateRepository[testUpdate](db, &Config{
			Validate: false,
		})

		err = repo2.Update(entity)
		assert.Nil(t, err)

		updated, err := repo.GetById(id)
		assert.Nil(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, entity.Email, updated.Email)
	})

	t.Run("Update error connection close", func(t *testing.T) {
		entity := createUpdateEntity(repo)
		assert.NotNil(t, entity)

		entity.Name = uuid.NewString()
		entity.Email = fmt.Sprintf("%v@mail.com", entity.Name)

		conn, err := repo.db.DB()
		assert.Nil(t, err)
		err = conn.Close()
		assert.Nil(t, err)

		err = repo.Update(entity)

		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "sql: database is closed")
	})
}
