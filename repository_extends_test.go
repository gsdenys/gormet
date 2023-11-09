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

type CustomRepo[T any] struct {
	Repository[T]
}

func (cr *CustomRepo[T]) teste() string {
	return "test"
}

func Test_Extends(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&ExtendModel{})

	bRepo, err := New[ExtendModel](db)
	repo := CustomRepo[ExtendModel]{
		*bRepo,
	}

	entity := &ExtendModel{
		Name: uuid.NewString(),
	}

	err = repo.Create(entity)
	assert.Nil(t, err)
	assert.NotZero(t, entity.ID)

	tst := repo.teste()
	assert.Equal(t, "test", tst)
}
