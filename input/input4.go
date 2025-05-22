package test

type Id uint
type Email string

type Card struct {
	ID    Id    `json:"id"`
	Email Email `json:"email"`
}
