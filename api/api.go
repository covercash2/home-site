package api

import (
	"github.com/gorilla/schema"
	"net/http"
)

var formDecoder = schema.NewDecoder()
var me Person

func InitAPI() {
	formDecoder.IgnoreUnknownKeys(true)
	me = Person{
		EmailAddress: "chris@covercash.biz",
		Name:         "Chris Overcash",
		PhoneNumber:  "5015100946",
	}
}

type emailForm struct {
	EmailAddress string `schema:"email"`
	Name         string `schema:"name"`
	PhoneNumber  string `schema:"phone"`
	Message      string `schema:"message"`
	Subject      string `schema:"subject"`
}

// ParseEmailForm takes a reference to a request,
// verifies it, and parses it into an Email object
func ParseEmailForm(request *http.Request) (Email, error) {
	var email Email
	var form emailForm
	var err error

	err = request.ParseForm()
	if err != nil {
		return email, err
	}

	values := request.Form
	err = formDecoder.Decode(&form, values)
	if err != nil {
		return email, err
	}

	return form.toEmail(me), nil
}

func (form *emailForm) toEmail(recipient Person) Email {
	return Email{
		Message:   form.Message,
		Recipient: recipient,
		Sender:    form.getSender(),
		Subject:   form.Subject,
	}
}

// TODO possibly search from database
func (form *emailForm) getSender() Person {
	return Person{
		EmailAddress: form.EmailAddress,
		Name:         form.Name,
		PhoneNumber:  form.PhoneNumber,
	}
}
