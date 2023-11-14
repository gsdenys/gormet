package gormet

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRepository_Search(t *testing.T) {
	type testSearch struct {
		gorm.Model
		Email  string `json:"email" gorm:"unique,not null;default:null"`
		Name   string `json:"name" gorm:"unique;not null;default:null"`
		Filter string `json:"filter"`
	}

	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, db)

	db.AutoMigrate(&testSearch{})

	repo, err := New[testSearch](db)
	assert.Nil(t, err)
	assert.NotNil(t, repo)

	createMany := func(numReg uint, group string) []*testSearch {
		var elements []*testSearch

		for n := 0; n < int(numReg); n++ {
			name := uuid.NewString()
			el := &testSearch{
				Name:   name,
				Email:  fmt.Sprintf("%s@email.com", name),
				Filter: group,
			}

			err := repo.Create(el)
			assert.Nil(t, err)

			elements = append(elements, el)

		}

		return elements
	}

	t.Run("Get 100 paginated elements", func(t *testing.T) {
		group := uuid.NewString()
		createdElements := createMany(100, group)

		repo.PageSize = 10
		assert.Equal(t, uint(10), repo.PageSize)

		resp, err := repo.Search(0, "filter = ?", group)

		assert.Nil(t, err)
		assert.NotNil(t, resp)

		assert.Equal(t, len(createdElements), int(resp.Response.TotalCount))
		assert.Equal(t, len(createdElements)/10, int(resp.Response.TotalPages))
		assert.Equal(t, 0, resp.Response.Page)
	})

	t.Run("Request Error", func(t *testing.T) {
		group := uuid.NewString()
		// createdElements := createMany(100, group)

		repo.PageSize = 10
		assert.Equal(t, uint(10), repo.PageSize)

		resp, err := repo.Search(0, "x = ?", group)

		assert.NotNil(t, err)
		assert.Empty(t, resp.Response.Entities)
	})

}
