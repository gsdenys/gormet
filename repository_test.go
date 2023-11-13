package gormet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNew(t *testing.T) {
	type TestNew struct {
		gorm.Model
		SomeField string `json:"somefield" gorm:"unique;not null;default:null" validate:"required,min=3,max=50"`
	}

	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&TestNew{})

	t.Run("Pagination active", func(t *testing.T) {
		repo, err := New[TestNew](db)

		assert.Nil(t, err)
		assert.NotNil(t, repo)

		assert.True(t, repo.Paginate)
	})

	t.Run("Pagination Inactive", func(t *testing.T) {
		repo, err := New[TestNew](db)
		assert.Nil(t, err)
		assert.NotNil(t, repo)

		repo.Paginate = false

		assert.False(t, repo.Paginate)
	})

	t.Run("Struct with custom pk", func(t *testing.T) {
		type CustomPK struct {
			MyID string `json:"myid" gorm:"primaryKey;autoIncrement:false"`
		}

		db.AutoMigrate(CustomPK{})

		repo, err := New[CustomPK](db)

		assert.Nil(t, err)
		assert.NotNil(t, repo)
		assert.Equal(t, "my_id", repo.pkName)
	})

	t.Run("Struct without pk", func(t *testing.T) {
		type NoPK struct {
			MyID string `json:"myid" gorm:"autoIncrement:false"`
		}

		db.AutoMigrate(NoPK{})

		repo, err := New[NoPK](db)

		assert.NotNil(t, err)
		assert.Nil(t, repo)
		assert.Equal(t, "impossible to retrieve primary key: no primary key found", err.Error())
	})
}

func Test_getPrimaryKeyFieldName(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	assert.Nil(t, err)

	t.Run("Default Primary Key", func(t *testing.T) {
		type DefaultPK struct {
			gorm.Model
			Name string `json:"name" gorm:"unique;not null;default:null"`
		}

		pk, err := getPrimaryKeyFieldName(db, DefaultPK{})
		assert.Nil(t, err)
		assert.NotNil(t, pk)
		assert.Equal(t, "id", pk)
	})

	t.Run("Custom Primary Key", func(t *testing.T) {
		type CustomPK struct {
			Email string `json:"email" gorm:"primaryKey"`
			Name  string `json:"name" gorm:"not null;default:null"`
		}

		pk, err := getPrimaryKeyFieldName(db, CustomPK{})
		assert.Nil(t, err)
		assert.NotNil(t, pk)
		assert.Equal(t, "email", pk)
	})

	t.Run("No Primary Key", func(t *testing.T) {
		type NoPk struct {
			Email string `json:"email" gorm:"unique;not null;default:null"`
			Name  string `json:"name" gorm:"not null;default:null"`
		}

		pk, err := getPrimaryKeyFieldName(db, NoPk{})
		assert.NotNil(t, err)
		assert.Equal(t, "", pk)
	})

	t.Run("Invalid struct", func(t *testing.T) {

		pk, err := getPrimaryKeyFieldName(db, "Just a parser error string")
		assert.NotNil(t, err)
		assert.Equal(t, "", pk)
	})
}
