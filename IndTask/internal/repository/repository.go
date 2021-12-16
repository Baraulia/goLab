package repository

type AppUser interface {
}

type AppBook interface {
}

type AppMove interface {
}

type Repository struct {
	AppUser
	AppBook
	AppMove
}

func NewRepository() *Repository {
	return &Repository{}
}
