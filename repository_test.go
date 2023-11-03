package gormet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type JustTest struct {
	gorm.Model
	SomeField string `json:"somefield" gorm:"unique;not null;default:null" validate:"required,min=3,max=50"`
}

func TestNew(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	db.AutoMigrate(&JustTest{})

	assert.Nil(t, err)
	assert.NotNil(t, db)

	repo := New[JustTest](db)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.db)
}

func TestCreateRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	db.AutoMigrate(&JustTest{})

	assert.Nil(t, err)
	assert.NotNil(t, db)

	repo := CreateRepository[JustTest](db, nil)
	assert.Equal(t, repo.config, DefaultConfig())

	repo = CreateRepository[JustTest](db, &Config{})
	assert.False(t, repo.config.Validate)
}
