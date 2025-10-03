package user

type User struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	// PassHash
	Password   string `json:"password"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

func NewUser(
	login,
	password,
	name,
	surname,
	patronymic,
	phone,
	email string,
) *User {
	return &User{
		Login:      login,
		Password:   password,
		Name:       name,
		Surname:    surname,
		Patronymic: patronymic,
		Phone:      phone,
		Email:      email,
	}
}

// Actions with user
type UserRepo interface {
	CreateUser(
		login string,
		password string,
		name string,
		surname string,
		patronymic string,
		phone string,
		email string,
	) (*User, error)
	GetUserByID(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
}
