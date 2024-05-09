package data

type Repository[T any] interface {
	GetAll() ([]T, error)
	Get(int) (T, error)
	Create(T) (T, error)
	Update(T) (T, error)
	Delete(int) error
}
