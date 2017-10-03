package api

import (
	"net/url"
	"testing"
)

var testForm url.Values

var testEmail = "chug.Mudfurd@mudfurd.co.uk"
var testName = "Chug Mudfurd"
var testPhoneNumber = "5555555555"
var testMessage = `a long and drawn out message
with some weird 'charachters' maybe 3-4? something
Like THAT. just dom #'s and (letters) to test
`
var testSubject = "testing"

func TestAPI(t *testing.T) {
	generateEmailForm()
	t.Run("parse email form", testParseEmailForm)
}

func generateEmailForm() {
	testForm = url.Values{}
	testForm.Set("email", testEmail)
	testForm.Set("name", testName)
	testForm.Set("phone", testPhoneNumber)
	testForm.Set("message", testMessage)
	testForm.Set("subject", testSubject)
}

func testParseEmailForm(testing *testing.T) {
	// generate mocks
	// verify the contents of the response
	// check Email conversion

	testing.Fail()
}
