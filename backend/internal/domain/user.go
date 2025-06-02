package domain

// User — здесь просто “обёртка” для хранения учётных данных администратора.
type User struct {
	Username string
	Password string
}
