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
	SomeField string `json:"somefield" gorm:"unique;not null;default:null" validate:"required,min=3,max=50"`
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

func TestRepository_Create(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&JustTest{})

	tests := []struct {
		name    string
		entity  *JustTest
		repo    *Repository[JustTest]
		wantErr string
	}{
		{
			name: "Valid JustTest Entity",
			entity: &JustTest{
				SomeField: fmt.Sprintf("test-%s", uuid.NewString()),
			},
			repo: New[JustTest](db),
		},
		{
			name:    "Empty JustTest Entity",
			entity:  &JustTest{},
			wantErr: "Key: 'JustTest.SomeField' Error:Field validation for 'SomeField' failed on the 'required' tag",
			repo:    New[JustTest](db),
		},
		{
			name:    "Nil Entity",
			entity:  nil,
			wantErr: "the entity cannot be nil",
			repo:    New[JustTest](db),
		},
		{
			name:    "Empty JustTest Entity without validation",
			entity:  &JustTest{},
			wantErr: "NOT NULL constraint failed: just_tests.some_field",
			repo:    CreateRepository[JustTest](db, &Config{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.repo.Create(tt.entity)

			if err == nil {
				assert.NotEqual(t, tt.entity.ID, 0)
			} else {
				assert.Equal(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestRepository_GetById(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&JustTest{})

	repo := New[JustTest](db)

	t.Run("Find New Entity", func(t *testing.T) {
		justTest := &JustTest{
			SomeField: uuid.NewString(),
		}

		if err := repo.Create(justTest); err != nil {
			t.Errorf("Create element error")
		}

		got, err := repo.GetById(justTest.ID)

		assert.Nil(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, got.ID, justTest.ID)
		assert.Equal(t, got.SomeField, justTest.SomeField)
	})

	t.Run("Element not found", func(t *testing.T) {
		got, err := repo.GetById(999999)

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "entity not found")
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

	db.AutoMigrate(&JustTest{})

	repo := New[JustTest](db)

	t.Run("Find new entity by field", func(t *testing.T) {
		justTest := &JustTest{
			SomeField: uuid.NewString(),
		}

		if err := repo.Create(justTest); err != nil {
			t.Errorf("Create element error")
		}

		got, err := repo.Get(JustTest{SomeField: justTest.SomeField})

		assert.Nil(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, got.ID, justTest.ID)
		assert.Equal(t, got.SomeField, justTest.SomeField)
	})

	t.Run("Find new entity by field with fail", func(t *testing.T) {
		where := JustTest{
			SomeField: "Isto é apenas um teste",
		}

		got, err := repo.Get(where)

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "record not found")
	})
}
