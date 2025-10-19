//////////////////// [ Day 5 ] ////////////////////

package user

type User struct {
	Name string
}

type IRepository interface {
	Get(id string) (User, error)
	Create(user User) error
}

type Service struct {
	Repository IRepository
}

func NewService(repo IRepository) *Service {
	return &Service{
		Repository: repo,
	}
}

func (service Service) Get(id string) (User, error) {
	return service.Repository.Get(id)
}

func (service Service) Create(user User) error {
	return service.Repository.Create(user)
}

//////////////////// [ Day 5 ] ////////////////////
