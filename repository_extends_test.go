package gormet

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ExtendModel struct {
	*gorm.Model
	Name string `json:"name" gorm:"unique;not null;default:null"`
}

type CustomStructExtendedRepo[T any] struct {
	Repository[T]
}

func (cs *CustomStructExtendedRepo[T]) SayHello() string {
	return "Hello Struct Extention"
}

func Test_Extended(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&ExtendModel{})

	t.Run("Struct Extention", func(t *testing.T) {
		r, _ := New[ExtendModel](db)
		repo := &CustomStructExtendedRepo[ExtendModel]{
			Repository: *r,
		}

		entity := &ExtendModel{
			Name: uuid.NewString(),
		}

		err := repo.Create(entity)

		assert.Nil(t, err)
		assert.NotZero(t, entity.ID)

		assert.Equal(t, "Hello Struct Extention", repo.SayHello())
	})
}
