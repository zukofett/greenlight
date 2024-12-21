package mocks

import (
	"context"
	"time"

	"github.com/zukofett/greenlight/internal/data"
)

var mockMetadata = data.Metadata{
	CurrentPage:  1,
	PageSize:     25,
	FirstPage:    1,
	LastPage:     1,
	TotalRecords: 1,
}

var mockMovie = data.Movie{
	ID:        1,
	CreatedAt: time.Now(),
	Title:     "Moana",
	Year:      2016,
	Runtime:   107,
	Genres:    []string{"animation", "adventure"},
	Version:   1,
}

type MovieModel struct{}

func (m MovieModel) Insert(_ context.Context, _ *data.Movie) error {
	return nil
}

func (m MovieModel) Get(_ context.Context, id int64) (*data.Movie, error) {
	switch id {
	case 1:
		return &mockMovie, nil
	default:
		return nil, data.ErrRecordNotFound
	}
}

func (m MovieModel) Update(_ context.Context, movie *data.Movie) error {
	switch movie.ID {
	case 1:
		return nil
	default:
		return data.ErrRecordNotFound
	}

}

func (m MovieModel) Delete(_ context.Context, id int64) error {
	switch id {
	case 1:
		return nil
	default:
		return data.ErrRecordNotFound
	}
}

func (m MovieModel) GetAll(_ context.Context, title string, genres []string, filters data.Filters) ([]*data.Movie, data.Metadata, error) {
	return []*data.Movie{&mockMovie}, mockMetadata, nil
}
