package gormet

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testGet struct {
	gorm.Model
	Name  string `json:"name" gorm:"unique;not null;default:null" validate:"required,min=3,max=50"`
	Email string `json:"email" gorm:"unique,not null;default:null" validate:"required,email"`
}

func createGetEntity(repo *Repository[testGet]) *testGet {
	field := uuid.NewString()
	entity := &testGet{
		Name:  field,
		Email: fmt.Sprintf("%v@mail.com", field),
	}

	if err := repo.Create(entity); err != nil {
		return nil
	}

	return entity
}

func TestRepository_GetById(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&testGet{})

	repo := New[testGet](db)

	t.Run("Get new entity by ID", func(t *testing.T) {
		entity := createGetEntity(repo)

		got, err := repo.GetById(entity.ID)

		assert.Nil(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, got.ID, entity.ID)
		assert.Equal(t, got.Name, entity.Name)
		assert.Equal(t, got.Email, entity.Email)
	})

	t.Run("Get by id error element not found", func(t *testing.T) {
		got, err := repo.GetById(99999999999)

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "record not found")
	})

	t.Run("Get by id error ", func(t *testing.T) {
		got, err := repo.GetById("this_is_not_an_id")

		assert.Nil(t, got)
		assert.NotNil(t, err)
	})

	t.Run("Nil id", func(t *testing.T) {
		got, err := repo.GetById(nil)

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "the id cannot be nil")
	})
}

func TestRepository_Get(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&testGet{})

	repo := New[testGet](db)

	t.Run("Get new entity", func(t *testing.T) {
		entity := createGetEntity(repo)

		got, err := repo.Get(*entity)

		assert.Nil(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, got.ID, entity.ID)
		assert.Equal(t, got.Name, entity.Name)
		assert.Equal(t, got.Email, entity.Email)
	})

	t.Run("Get first record", func(t *testing.T) {
		got, err := repo.Get(testGet{})

		assert.Nil(t, err)
		assert.NotNil(t, got)

		var first testGet
		repo.db.First(&first)

		assert.Equal(t, first.ID, got.ID)
		assert.Equal(t, first.Name, got.Name)
		assert.Equal(t, first.Email, got.Email)
	})

	t.Run("Get with closed connection", func(t *testing.T) {
		sql, _ := db.DB()
		sql.Close()

		got, err := repo.Get(testGet{})

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "sql: database is closed")
	})
}
