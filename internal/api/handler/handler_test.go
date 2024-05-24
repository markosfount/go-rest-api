package handler

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"rest_api/internal/api/model"
	"rest_api/internal/api/service"
	"rest_api/internal/api/tmdb"
	"strings"
	"testing"
)

type mockMovieRepository struct {
	mock.Mock
}

func (r *mockMovieRepository) GetAll() ([]*model.Movie, error) {
	args := r.Called()
	return args.Get(0).([]*model.Movie), args.Error(1)
}

func (r *mockMovieRepository) Get(movieId int) (*model.Movie, error) {
	args := r.Called(movieId)
	return args.Get(0).(*model.Movie), args.Error(1)
}

func (r *mockMovieRepository) Create(movie *model.Movie) (*model.Movie, error) {
	args := r.Called(movie)
	return args.Get(0).(*model.Movie), args.Error(1)
}

func (r *mockMovieRepository) Update(movie *model.Movie) (*model.Movie, error) {
	args := r.Called(movie)
	return args.Get(0).(*model.Movie), args.Error(1)
}

func (r *mockMovieRepository) Delete(movieId int) error {
	args := r.Called(movieId)
	return args.Error(0)
}

type mockPublisher struct {
	mock.Mock
}

func (p *mockPublisher) Publish(msg string) error {
	args := p.Called(msg)
	return args.Error(0)
}

func (p *mockPublisher) Configure(_ string) {}

func TestHandler_GetMovie(t *testing.T) {
	w := httptest.NewRecorder()

	repository := new(mockMovieRepository)
	// TODO change anytnhing
	repository.On("Get", mock.Anything).Return(&model.Movie{1, "foo", "bar"}, nil)

	h := Handler{
		UserRepository: nil,
		MovieService:   service.NewMovieService(repository),
		TmdbService:    nil,
		Publisher:      nil,
	}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:3000/movies/", nil)
	req = mux.SetURLVars(req, map[string]string{"movieId": "1"})
	require.NoError(t, err)

	h.GetMovie(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	bytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, `{"id":1,"title":"foo","overview":"bar"}`, string(bytes))
}

func TestHandler_CreateMovie(t *testing.T) {
	tmdbServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":[{"id":1,"original_title":"The bear","overview":"bear","runtime":123}],"total_results": 1}`))
	}))
	defer tmdbServer.Close()

	w := httptest.NewRecorder()

	mockRepository := new(mockMovieRepository)
	mockRepository.On("Create", mock.Anything).Return(&model.Movie{1, "The bear", "bear"}, nil)

	publisher := new(mockPublisher)
	publisher.On("Publish", mock.Anything).Return(nil)

	h := Handler{
		UserRepository: nil,
		MovieService:   service.NewMovieService(mockRepository),
		TmdbService:    tmdb.NewService(tmdbServer.URL),
		Publisher:      publisher,
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:3000/movies/", strings.NewReader(`{"id":45,"title":"The bear"}`))
	require.NoError(t, err)

	h.AddMovie(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	bytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, `{"id":1,"title":"The bear","overview":"bear"}`, string(bytes))
}
