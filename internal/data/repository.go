package data

type Repository[T any] interface {
	GetAll() ([]T, error)
	Get(string) (T, error)
	Create(T) (T, error)
	Update(T) (T, error)
	Delete(string) error
}
