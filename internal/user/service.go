package user

import "log"

type (
	Filters struct {
		Firstname string
		Lastname  string
	}

	Service interface {
		Create(firstname, lastname, email, phone string) (*User, error)
		Get(id string) (*User, error)
		GetAll(filters Filters, offset, limit int) ([]User, error)
		Delete(id string) (*User, error)
		Update(id string, firstname, lastname, email, phone *string) (*User, error)
		Count(Filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func InitService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(firstname, lastname, email, phone string) (*User, error) {
	return s.repo.Create(&User{
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Phone:     phone,
	})
}

func (s service) GetAll(filters Filters, offset, limit int) ([]User, error) {
	return s.repo.FindAll(filters, offset, limit)
}

func (s service) Get(id string) (*User, error) {
	return s.repo.FindOne(id)
}

func (s service) Delete(id string) (*User, error) {
	return s.repo.Delete(id)
}

func (s service) Update(id string, firstname, lastname, email, phone *string) (*User, error) {
	return s.repo.Update(id, firstname, lastname, email, phone)
}

func (s service) Count(filters Filters) (int, error) {
	return s.repo.Count(filters)
}
