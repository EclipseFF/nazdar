package internal

type Password struct {
	Plaintext *string
	Hash *string
}

type Admin struct {
	Id       *int
	Username *string
	Password *Password
	Permissions *string
}

