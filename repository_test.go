package gormet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNew(t *testing.T) {
	type JustTest struct {
		gorm.Model
		SomeField string `json:"somefield" gorm:"unique;not null;default:null" validate:"required,min=3,max=50"`
	}

	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&JustTest{})

	t.Run("Default config", func(t *testing.T) {
		repo, err := New[JustTest](db)

		assert.Nil(t, err)
		assert.NotNil(t, repo)

		assert.Equal(t, DefaultConfig(), repo.config)
	})

	t.Run("Custom config", func(t *testing.T) {
		repo, err := New[JustTest](db)

		assert.Nil(t, err)
		assert.NotNil(t, repo)

		repo.Validate(false)

		assert.NotEqual(t, DefaultConfig(), repo.config)
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

func TestRepository_Paginate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	assert.Nil(t, err)

	t.Run("Set paginate false", func(t *testing.T) {
		type SetPaginate struct {
			gorm.Model
			Name string `json:"name" gorm:"autoIncrement:false"`
		}

		db.AutoMigrate(SetPaginate{})

		repo, err := New[SetPaginate](db)
		assert.Nil(t, err)

		repo.Paginate(false)
		assert.False(t, repo.config.Paginate)
	})
}

func TestRepository_Validate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	assert.Nil(t, err)

	t.Run("Set Validate false", func(t *testing.T) {
		type SetPaginate struct {
			gorm.Model
			Name string `json:"name" gorm:"autoIncrement:false"`
		}

		db.AutoMigrate(SetPaginate{})

		repo, err := New[SetPaginate](db)
		assert.Nil(t, err)

		repo.Validate(false)
		assert.False(t, repo.config.Validate)
	})
}
