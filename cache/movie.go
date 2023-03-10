package cache

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type Movie struct {
	Id          string `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type MovieService interface {
	GetMovie(id string) (*Movie, error)
	GetMovies() ([]*Movie, error)
	CreateMovie(movie *Movie) (*Movie, error)
	UpdateMovie(movie *Movie) (*Movie, error)
	DeleteMovie(id string) error
}

func (cache redisCache) CreateMovie(movie *Movie) (*Movie, error) {
	c := cache.getClient()
	movie.Id = uuid.New().String()
	json, err := json.Marshal(movie)
	if err != nil {
		return nil, err
	}
	c.HSet("movies", movie.Id, json)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (cache redisCache) GetMovie(id string) (*Movie, error) {
	c := cache.getClient()
	val, err := c.HGet("movies", id).Result()

	if err != nil {
		return nil, err
	}
	movie := &Movie{}
	err = json.Unmarshal([]byte(val), movie)

	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (cache redisCache) GetMovies() ([]*Movie, error) {
	c := cache.getClient()
	movies := []*Movie{}
	val, err := c.HGetAll("movies").Result()
	if err != nil {
		return nil, err
	}
	for _, item := range val {
		movie := &Movie{}
		err := json.Unmarshal([]byte(item), movie)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (cache redisCache) UpdateMovie(movie *Movie) (*Movie, error) {
	c := cache.getClient()
	json, err := json.Marshal(&movie)
	if err != nil {
		return nil, err
	}
	c.HSet("movies", movie.Id, json)
	if err != nil {
		return nil, err
	}
	return movie, nil
}
func (cache redisCache) DeleteMovie(id string) error {
	c := cache.getClient()
	numDeleted, err := c.HDel("movies", id).Result()
	if numDeleted == 0 {
		return errors.New("movie to delete not found")
	}
	if err != nil {
		return err
	}
	return nil
}
