package entities

type User struct {
	Login       string `db:"login"`
	Email       string `db:"email"`
	Password    string `db:"password_hash"`
	PhoneNumber string `db:"phone_number"`
}
