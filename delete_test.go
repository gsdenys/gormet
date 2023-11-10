package gormet

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRepository_DeleteById(t *testing.T) {

	type testDelete struct {
		gorm.Model
		Email string `json:"email" gorm:"unique,not null;default:null"`
		Name  string `json:"name" gorm:"unique;not null;default:null"`
	}

	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, db)

	db.AutoMigrate(&testDelete{})

	repo, err := New[testDelete](db)
	assert.Nil(t, err)

	createRegister := func() *testDelete {
		var field = uuid.NewString()
		entity := &testDelete{
			Name:  fmt.Sprintf("%v", field),
			Email: fmt.Sprintf("%v@mail.com", field),
		}

		repo.Create(entity)

		return entity
	}

	t.Run("Updated Successfull", func(t *testing.T) {
		entity := createRegister()
		id := entity.ID

		err := repo.DeleteById(id)
		assert.Nil(t, err)

		var got *testDelete = &testDelete{}
		tx := db.First(got, "id = ?", id)

		assert.NotNil(t, tx)
		assert.NotNil(t, tx.Error)
		assert.Equal(t, "record not found", tx.Error.Error())
		assert.Equal(t, uint(0), got.ID)
	})

	t.Run("Nil entity", func(t *testing.T) {
		err := repo.DeleteById(nil)

		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "the ID should not be nil")
	})

	t.Run("Entity with sql injection", func(t *testing.T) {
		entity := createRegister()

		id := fmt.Sprintf(`%d"; update test_deletes set deleted at '1970-01-01 00:00:00 --`, entity.ID)

		err = repo.DeleteById(id)
		assert.NotNil(t, err)
		assert.Equal(t, "no register found", err.Error())
	})

	t.Run("Connection Closed", func(t *testing.T) {
		entity := createRegister()

		sqlDB, _ := db.DB()
		err := sqlDB.Close()
		assert.Nil(t, err)

		err = repo.DeleteById(entity.ID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "sql: database is closed")
	})
}
