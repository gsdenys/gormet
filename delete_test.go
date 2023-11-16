package gormet

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type testDelete struct {
	gorm.Model
	Email string `json:"email" gorm:"unique,not null;default:null"`
	Name  string `json:"name" gorm:"unique;not null;default:null"`
}

func TestRepository_DeleteById(t *testing.T) {
	db := getGormConnection(t, &testDelete{})

	repo, err := New[testDelete](db)
	assert.Nil(t, err)

	t.Run("Delete entity successfully", func(t *testing.T) {
		entity := &testDelete{
			Name:  uuid.NewString(),
			Email: fmt.Sprintf("%s@mail.com", uuid.NewString()),
		}
		err := repo.Create(entity)
		assert.Nil(t, err)

		id := entity.ID

		err = repo.DeleteById(id)
		assert.Nil(t, err)

		var got *testDelete = &testDelete{}
		tx := db.First(got, "id = ?", id)

		assert.NotNil(t, tx)
		assert.NotNil(t, tx.Error)
		assert.Equal(t, "record not found", tx.Error.Error())
		assert.Equal(t, uint(0), got.ID)
	})

	t.Run("Error send nil entity", func(t *testing.T) {
		err := repo.DeleteById(nil)

		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "the ID should not be nil")
	})

	t.Run("Entity with sql injection", func(t *testing.T) {
		entity := &testDelete{
			Name:  uuid.NewString(),
			Email: fmt.Sprintf("%s@mail.com", uuid.NewString()),
		}
		err := repo.Create(entity)
		assert.Nil(t, err)

		id := fmt.Sprintf("%d", entity.ID) + `"; update test_deletes set deleted_at '1970-01-01 00:00:00 --`

		err = repo.DeleteById(id)
		assert.NotNil(t, err)
		assert.Equal(t, "no register found", err.Error())
	})

	t.Run("Connection Closed", func(t *testing.T) {
		sqlDB, _ := db.DB()
		err := sqlDB.Close()
		assert.Nil(t, err)

		err = repo.DeleteById(1)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "sql: database is closed")
	})
}

func TestRepository_Delete(t *testing.T) {
	db := getGormConnection(t, &testDelete{})

	repo, err := New[testDelete](db)
	assert.Nil(t, err)

	t.Run("Delete entity successfully", func(t *testing.T) {
		entity := &testDelete{
			Name:  uuid.NewString(),
			Email: fmt.Sprintf("%s@mail.com", uuid.NewString()),
		}
		err := repo.Create(entity)
		assert.Nil(t, err)

		id := entity.ID

		err = repo.Delete(entity)
		assert.Nil(t, err)

		var got *testDelete = &testDelete{}
		tx := db.First(got, "id = ?", id)

		assert.NotNil(t, tx)
		assert.NotNil(t, tx.Error)
		assert.Equal(t, "record not found", tx.Error.Error())
		assert.Equal(t, uint(0), got.ID)
	})

	t.Run("Error send nil entity", func(t *testing.T) {
		err := repo.Delete(nil)

		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "the entity should not be nil")
	})

	t.Run("Entity with sql injection", func(t *testing.T) {
		entity := &testDelete{
			Name:  fmt.Sprintf("sql-%s", uuid.NewString()),
			Email: fmt.Sprintf("%s@sql.com", uuid.NewString()),
		}
		err := repo.Create(entity)
		assert.Nil(t, err)

		name := entity.Name + `"; update test_deletes set deleted at '1970-01-01 00:00:00 --`
		entity.Name = name

		err = repo.Delete(entity)
		assert.Nil(t, err)
	})

	t.Run("Entity with sql injection in string ID", func(t *testing.T) {
		type testDeletePkString struct {
			Email string `json:"email" gorm:"primaryKey;default:null"`
			Name  string `json:"name" gorm:"default:null"`
		}

		db.AutoMigrate(&testDeletePkString{})

		repo, err := New[testDeletePkString](db)
		assert.Nil(t, err)

		entity := &testDeletePkString{}
		name := uuid.NewString() + `"; update test_deletes set deleted at '1970-01-01 00:00:00 --`
		entity.Email = name

		err = repo.Delete(entity)
		assert.NotNil(t, err)
		assert.Equal(t, "no register found", err.Error())
	})

	t.Run("Error delete without where", func(t *testing.T) {
		repo, err := New[testDelete](db)
		assert.Nil(t, err)

		entity := &testDelete{}

		err = repo.Delete(entity)
		assert.NotNil(t, err)
		assert.Equal(t, "WHERE conditions required", err.Error())
	})

	t.Run("Connection Closed", func(t *testing.T) {
		sqlDB, _ := db.DB()
		err := sqlDB.Close()
		assert.Nil(t, err)

		err = repo.DeleteById(2)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "sql: database is closed")
	})
}
