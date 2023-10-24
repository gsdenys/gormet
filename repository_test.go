package gormet

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type JustTest struct {
	gorm.Model
	SomeField string `json:"somefield" gorm:"unique;not null" validate:"required,min=3,max=50"`
}

func TestNew(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	db.AutoMigrate(&JustTest{})

	assert.Nil(t, err)
	assert.NotNil(t, db)

	repo := New[JustTest](db)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.DB)
}

func TestCreateRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	db.AutoMigrate(&JustTest{})

	assert.Nil(t, err)
	assert.NotNil(t, db)

	repo := CreateRepository[JustTest](db, nil)
	assert.Equal(t, repo.Config, DefaultConfig())

	repo = CreateRepository[JustTest](db, &Config{})
	assert.False(t, repo.Config.Validate)
}

// func TestRepository_Create(t *testing.T) {
// 	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
// 	assert.Nil(t, err)

// 	db.AutoMigrate(&JustTest{})

// 	repo := New[JustTest](db)

// 	entity := &JustTest{
// 		SomeField: uuid.NewString(),
// 	}

// 	err = repo.Create(entity)

// 	assert.Nil(t, err)
// 	assert.NotEqual(t, entity.ID, 0)
// }

func TestRepository_Create(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&JustTest{})

	repo := New[JustTest](db)

	tests := []struct {
		name    string
		entity  *JustTest
		wantErr string
	}{
		{
			name: "Valid JustTest Entity",
			entity: &JustTest{
				SomeField: fmt.Sprintf("test-%s", uuid.NewString()),
			},
		},
		{
			name:    "Empty JustTest Entity",
			entity:  &JustTest{},
			wantErr: "Key: 'JustTest.SomeField' Error:Field validation for 'SomeField' failed on the 'required' tag",
		},
		{
			name:    "Nil Entity",
			entity:  nil,
			wantErr: "the entity cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.entity)

			if err == nil {
				assert.NotEqual(t, tt.entity.ID, 0)
			} else {
				assert.Equal(t, err.Error(), tt.wantErr)
			}

		})
	}
}
