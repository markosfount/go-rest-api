package service

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"reflect"
	"rest_api/internal/api/model"
	"rest_api/internal/data"
	"testing"
)

type MockRepository struct {
	mock.Mock
}

func (r *MockRepository) GetAll() ([]*model.Movie, error) {
	args := r.Called()
	return args.Get(0).([]*model.Movie), args.Error(1)
}

func (r *MockRepository) Get(movieId string) (*model.Movie, error) {
	args := r.Called(movieId)
	arg1 := args.Get(0)
	if arg1 != nil {
		return arg1.(*model.Movie), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *MockRepository) Create(movie *model.Movie) (*model.Movie, error) {
	args := r.Called(movie)
	arg1 := args.Get(0)
	if arg1 != nil {
		return arg1.(*model.Movie), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *MockRepository) Update(movie *model.Movie) (*model.Movie, error) {
	args := r.Called(movie)
	arg1 := args.Get(0)
	if arg1 != nil {
		return arg1.(*model.Movie), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *MockRepository) Delete(movieId string) error {
	args := r.Called(movieId)
	return args.Error(0)
}

func TestMovieService_GetAll(t *testing.T) {
	randErr := errors.New("random")
	mockRepository := MockRepository{}
	s := &MovieService{
		MovieRepository: &mockRepository,
	}
	tests := []struct {
		name     string
		want     []*model.Movie
		wantErr  error
		mockFunc func(r *MockRepository) *mock.Call
	}{
		{
			"success",
			[]*model.Movie{{MovieId: "1"}},
			nil,
			func(r *MockRepository) *mock.Call {
				return r.On("GetAll", mock.Anything).Return([]*model.Movie{{MovieId: "1"}}, nil)
			},
		},
		{
			"other error",
			[]*model.Movie{},
			randErr,
			func(r *MockRepository) *mock.Call {
				return r.On("GetAll", mock.Anything).Return([]*model.Movie{}, randErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCall := tt.mockFunc(&mockRepository)
			defer mockCall.Unset()
			got, err := s.GetAll()
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				if err != nil {
					t.Errorf("GetAll() error = %v, wantErr is nil", err)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMovieService_Get(t *testing.T) {
	randErr := errors.New("random")
	mockRepository := MockRepository{}
	s := &MovieService{
		MovieRepository: &mockRepository,
	}
	tests := []struct {
		name     string
		input    string
		want     *model.Movie
		wantErr  error
		mockFunc func(r *MockRepository) *mock.Call
	}{
		{
			"success",
			"1",
			&model.Movie{MovieId: "1"},
			nil,
			func(r *MockRepository) *mock.Call {
				return r.On("Get", mock.Anything).Return(&model.Movie{MovieId: "1"}, nil)
			},
		},
		{
			"not found",
			"1",
			nil,
			model.NotFoundError{},
			func(r *MockRepository) *mock.Call {
				return r.On("Get", mock.Anything).Return(nil, data.ErrRecordNotFound)
			},
		},
		{
			"other error",
			"1",
			nil,
			randErr,
			func(r *MockRepository) *mock.Call {
				return r.On("Get", mock.Anything).Return(nil, randErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCall := tt.mockFunc(&mockRepository)
			defer mockCall.Unset()
			got, err := s.Get(tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				if err != nil {
					t.Errorf("Get() error = %v, wantErr is nil", err)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMovieService_Create(t *testing.T) {
	randErr := errors.New("random")
	mockRepository := MockRepository{}
	s := &MovieService{
		MovieRepository: &mockRepository,
	}
	tests := []struct {
		name     string
		input    *model.Movie
		want     *model.Movie
		wantErr  error
		mockFunc func(r *MockRepository) *mock.Call
	}{
		{
			"success",
			&model.Movie{},
			&model.Movie{MovieId: "1"},
			nil,
			func(r *MockRepository) *mock.Call {
				return r.On("Create", mock.Anything).Return(&model.Movie{MovieId: "1"}, nil)
			},
		},
		{
			"conflict",
			&model.Movie{},
			nil,
			model.ConflictError{},
			func(r *MockRepository) *mock.Call {
				return r.On("Create", mock.Anything).Return(nil, data.ErrRecordExists)
			},
		},
		{
			"other error",
			&model.Movie{},
			nil,
			randErr,
			func(r *MockRepository) *mock.Call {
				return r.On("Create", mock.Anything).Return(nil, randErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCall := tt.mockFunc(&mockRepository)
			defer mockCall.Unset()
			got, err := s.Create(tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				if err != nil {
					t.Errorf("Create() error = %v, wantErr is nil", err)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMovieService_Update(t *testing.T) {
	randErr := errors.New("random")
	mockRepository := MockRepository{}
	s := &MovieService{
		MovieRepository: &mockRepository,
	}
	tests := []struct {
		name     string
		input    *model.Movie
		want     *model.Movie
		wantErr  error
		mockFunc func(r *MockRepository) *mock.Call
	}{
		{
			"success",
			&model.Movie{},
			&model.Movie{MovieId: "1"},
			nil,
			func(r *MockRepository) *mock.Call {
				return r.On("Update", mock.Anything).Return(&model.Movie{MovieId: "1"}, nil)
			},
		},
		{
			"not found",
			&model.Movie{},
			nil,
			model.NotFoundError{},
			func(r *MockRepository) *mock.Call {
				return r.On("Update", mock.Anything).Return(nil, data.ErrRecordNotFound)
			},
		},
		{
			"other error",
			&model.Movie{},
			nil,
			randErr,
			func(r *MockRepository) *mock.Call {
				return r.On("Update", mock.Anything).Return(nil, randErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCall := tt.mockFunc(&mockRepository)
			defer mockCall.Unset()
			got, err := s.Update(tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				if err != nil {
					t.Errorf("Update() error = %v, wantErr is nil", err)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMovieService_Delete(t *testing.T) {
	randErr := errors.New("random")
	mockRepository := MockRepository{}
	s := &MovieService{
		MovieRepository: &mockRepository,
	}
	tests := []struct {
		name     string
		input    string
		wantErr  error
		mockFunc func(r *MockRepository) *mock.Call
	}{
		{
			"success",
			"1",
			nil,
			func(r *MockRepository) *mock.Call {
				return r.On("Delete", mock.Anything).Return(nil)
			},
		},
		{
			"not found",
			"1",
			nil,
			func(r *MockRepository) *mock.Call {
				return r.On("Delete", mock.Anything).Return(data.ErrRecordNotFound)
			},
		},
		{
			"other error",
			"1",
			randErr,
			func(r *MockRepository) *mock.Call {
				return r.On("Delete", mock.Anything).Return(randErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCall := tt.mockFunc(&mockRepository)
			defer mockCall.Unset()
			err := s.Delete(tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				if err != nil {
					t.Errorf("Delete() error = %v, wantErr is nil", err)
				}
			}
		})
	}
}
