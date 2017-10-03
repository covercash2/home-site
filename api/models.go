package api

type Person struct {
	EmailAddress string `schema:"email"`
	Name         string `schema:"name"`
	PhoneNumber  string `schema:"phone"`
}

type Email struct {
	Sender    Person
	Recipient Person
	Subject   string
	Message   string
}
