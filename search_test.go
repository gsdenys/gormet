package gormet

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testSearch struct {
	gorm.Model
	Email  string `json:"email" gorm:"unique,not null;default:null"`
	Name   string `json:"name" gorm:"unique;not null;default:null"`
	Filter string `json:"filter"`
}

func createMany(repo *Repository[testSearch], numReg uint, group string) []*testSearch {
	var elements []*testSearch

	for n := 0; n < int(numReg); n++ {
		name := uuid.NewString()
		el := &testSearch{
			Name:   name,
			Email:  fmt.Sprintf("%s@email.com", name),
			Filter: group,
		}

		repo.Create(el)

		elements = append(elements, el)
	}

	return elements
}

func TestRepository_Search(t *testing.T) {

	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, db)

	db.AutoMigrate(&testSearch{})

	repo, err := New[testSearch](db)
	assert.Nil(t, err)
	assert.NotNil(t, repo)

	t.Run("Get first page", func(t *testing.T) {
		group := uuid.NewString()
		createdElements := createMany(repo, 100, group)

		repo.PageSize = 10
		assert.Equal(t, uint(10), repo.PageSize)

		resp, err := repo.Search(1, "filter = ?", group)

		assert.Nil(t, err)
		assert.NotNil(t, resp)

		assert.Equal(t, len(createdElements), int(resp.Response.TotalCount))
		assert.Equal(t, len(createdElements)/10, int(resp.Response.TotalPages))
		assert.Equal(t, uint(1), resp.Response.Page)
		assert.True(t, resp.Response.HasNextPage)
		assert.False(t, resp.Response.HasPrevPage)
	})

	t.Run("Second page", func(t *testing.T) {
		group := uuid.NewString()
		createdElements := createMany(repo, 100, group)

		repo.PageSize = 10
		assert.Equal(t, uint(10), repo.PageSize)

		resp, err := repo.Search(2, "filter = ?", group)

		assert.Nil(t, err)
		assert.NotNil(t, resp)

		assert.Equal(t, len(createdElements), int(resp.Response.TotalCount))
		assert.Equal(t, len(createdElements)/10, int(resp.Response.TotalPages))
		assert.Equal(t, uint(2), resp.Response.Page)
		assert.True(t, resp.Response.HasNextPage)
		assert.True(t, resp.Response.HasPrevPage)
	})

	t.Run("Last page", func(t *testing.T) {
		group := uuid.NewString()
		createdElements := createMany(repo, 100, group)

		repo.PageSize = 10
		assert.Equal(t, uint(10), repo.PageSize)

		resp, err := repo.Search(10, "filter = ?", group)

		assert.Nil(t, err)
		assert.NotNil(t, resp)

		assert.Equal(t, len(createdElements), int(resp.Response.TotalCount))
		assert.Equal(t, len(createdElements)/10, int(resp.Response.TotalPages))
		assert.Equal(t, uint(10), resp.Response.Page)
		assert.False(t, resp.Response.HasNextPage)
		assert.True(t, resp.Response.HasPrevPage)
	})

	t.Run("Request Error", func(t *testing.T) {
		group := uuid.NewString()

		repo.PageSize = 10
		assert.Equal(t, uint(10), repo.PageSize)

		resp, err := repo.Search(1, "this_field_generate_error = ?", group)

		assert.NotNil(t, err)
		assert.Empty(t, resp.Response.Entities)
	})

}

func Test_getOffset(t *testing.T) {
	type args struct {
		page     uint
		pageSize uint
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Page: 1, Page Size: 10",
			args: args{
				page:     1,
				pageSize: 10,
			},
			want: 0,
		},
		{
			name: "Page: 2, Page Size: 10",
			args: args{
				page:     2,
				pageSize: 10,
			},
			want: 10,
		},
		{
			name: "Page: 0, Page Size: 10",
			args: args{
				page:     0,
				pageSize: 10,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOffset(tt.args.page, tt.args.pageSize); got != tt.want {
				t.Errorf("defOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLimit(t *testing.T) {
	type args struct {
		pageSize uint
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Page Size 0",
			args: args{pageSize: 0},
			want: -1,
		},
		{
			name: "Page Size 1",
			args: args{pageSize: 1},
			want: 1,
		},
		{
			name: "Page Size 50",
			args: args{pageSize: 50},
			want: 50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLimit(tt.args.pageSize); got != tt.want {
				t.Errorf("getLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}
