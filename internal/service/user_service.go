package service

type IUserService interface {
    GetAllUsers() ([]string, error)
}

type UserService struct{}

func NewUserService() IUserService {
    return &UserService{}
}

func (s *UserService) GetAllUsers() ([]string, error) {
    return []string{"Alice", "Bob"}, nil
}
