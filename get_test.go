package gormet

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// testGet data structure to be used in test case
type testGet struct {
	gorm.Model
	Email string `json:"email" gorm:"unique,not null;default:null"`
	Name  string `json:"name" gorm:"unique;not null;default:null"`
}

// createTestGetRegister auxiliar function to help create entity to test case
func createTestGetRegister(repo *Repository[testGet]) *testGet {
	var field = uuid.NewString()
	entity := &testGet{
		Name:  fmt.Sprintf("%v", field),
		Email: fmt.Sprintf("%v@mail.com", field),
	}

	repo.Create(entity)

	return entity
}

func TestRepository_GetById(t *testing.T) {

	db := getGormConnection(t, &testGet{})

	repo, err := New[testGet](db)
	assert.Nil(t, err)

	t.Run("Get one register successfully", func(t *testing.T) {
		entity := createTestGetRegister(repo)
		id := entity.ID

		got, err := repo.GetById(id)

		assert.NotNil(t, got)
		assert.Nil(t, err)
		assert.Equal(t, id, got.ID)
	})

	t.Run("Error pass a nil id", func(t *testing.T) {
		got, err := repo.GetById(nil)

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, "the id should not be nil", err.Error())
	})

	t.Run("Error id with sql injection", func(t *testing.T) {
		entity := createTestGetRegister(repo)

		id := fmt.Sprintf("%d", entity.ID) + `"; update test_gets set deleted at '1970-01-01 00:00:00 --`

		got, err := repo.GetById(id)

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, "record not found", err.Error())
	})

	t.Run("Error connection closed", func(t *testing.T) {
		entity := createTestGetRegister(repo)

		sqlDB, _ := db.DB()
		err := sqlDB.Close()
		assert.Nil(t, err)

		got, err := repo.GetById(entity.ID)

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "sql: database is closed")
	})
}

func TestRepository_Get(t *testing.T) {

	db := getGormConnection(t, &testGet{})

	repo, err := New[testGet](db)
	assert.Nil(t, err)

	t.Run("Get one register successfully", func(t *testing.T) {
		entity := createTestGetRegister(repo)
		id := entity.ID

		got, err := repo.Get(testGet{
			Name:  entity.Name,
			Email: entity.Email,
		})

		assert.NotNil(t, got)
		assert.Nil(t, err)
		assert.Equal(t, id, got.ID)
	})

	t.Run("Entity with sql injection", func(t *testing.T) {
		entity := createTestGetRegister(repo)

		name := fmt.Sprintf("%d", entity.ID) + `"; update test_gets set deleted at '1970-01-01 00:00:00 --`
		entity.Name = name

		got, err := repo.Get(testGet{Name: name})

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, "record not found", err.Error())
	})

	t.Run("Error connection closed", func(t *testing.T) {
		entity := createTestGetRegister(repo)

		sqlDB, _ := db.DB()
		err := sqlDB.Close()
		assert.Nil(t, err)

		got, err := repo.Get(*entity)

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "sql: database is closed")
	})
}
