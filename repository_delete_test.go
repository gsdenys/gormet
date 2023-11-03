package gormet

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testDelete struct {
	gorm.Model
	Name  string `json:"name" gorm:"unique;not null;default:null" validate:"required,min=3,max=50"`
	Email string `json:"email" gorm:"unique,not null;default:null" validate:"required,email"`
}

func createDeleteEntity(repo *Repository[testDelete]) *testDelete {
	field := uuid.NewString()
	entity := &testDelete{
		Name:  field,
		Email: fmt.Sprintf("%v@mail.com", field),
	}

	if err := repo.Create(entity); err != nil {
		return nil
	}

	return entity
}

func TestRepository_DeleteById(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, db)

	db.AutoMigrate(&testDelete{})

	repo := New[testDelete](db)

	t.Run("Delete entity successfull", func(t *testing.T) {
		entity := createDeleteEntity(repo)
		assert.NotNil(t, entity)

		err := repo.DeleteById(entity.ID)
		assert.Nil(t, err)

		deleted, err := repo.GetById(entity.ID)
		assert.NotNil(t, err)
		assert.Nil(t, deleted)
		assert.Equal(t, err.Error(), "record not found")
	})

	t.Run("Delete error nil id", func(t *testing.T) {
		err := repo.DeleteById(nil)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "the id cannot be nil")
	})

	// t.Run("Delete error invalid id", func(t *testing.T) {
	// 	err := repo.DeleteById("some test")
	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, err.Error(), "")
	// })

	t.Run("my test", func(t *testing.T) {
		userInput := "10; update table test_deletes set deleted_at = null --"

		// Use parameterized queries to prevent SQL injection
		var entity testDelete
		repo.db.Debug().Where("id = ?", userInput).First(&entity)
	})

	// t.Run("Update invalid entity not validated", func(t *testing.T) {
	// 	entity := createUpdateEntity(repo)
	// 	assert.NotNil(t, entity)

	// 	id := entity.ID

	// 	entity.Name = uuid.NewString()
	// 	entity.Email = uuid.NewString() //fail in validation

	// 	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	// 	assert.Nil(t, err)
	// 	assert.NotNil(t, db)

	// 	repo2 := CreateRepository[testUpdate](db, &Config{
	// 		Validate: false,
	// 	})

	// 	err = repo2.Update(entity)
	// 	assert.Nil(t, err)

	// 	updated, err := repo.GetById(id)
	// 	assert.Nil(t, err)
	// 	assert.NotNil(t, updated)
	// 	assert.Equal(t, entity.Email, updated.Email)
	// })

	// t.Run("Update error connection close", func(t *testing.T) {
	// 	entity := createUpdateEntity(repo)
	// 	assert.NotNil(t, entity)

	// 	entity.Name = uuid.NewString()
	// 	entity.Email = fmt.Sprintf("%v@mail.com", entity.Name)

	// 	conn, err := repo.db.DB()
	// 	assert.Nil(t, err)
	// 	err = conn.Close()
	// 	assert.Nil(t, err)

	// 	err = repo.Update(entity)

	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, err.Error(), "sql: database is closed")
	// })
}
