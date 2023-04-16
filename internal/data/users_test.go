package data

import (
	"github.com/Skaifai/gophers-online-store/internal/validator"
	"testing"
	"time"
)

func TestValidateEmail(t *testing.T) {
	v := validator.New()

	emptyEmail := ""

	ValidateEmail(v, emptyEmail)
	expected := v.Valid()

	if expected != false {
		t.Errorf("ValidateProduct() returned unexpected value: got %v, expected %s", expected, "false")
	}
}

func TestTableDrivenValidateEmail(t *testing.T) {
	specialCharactersInFirstPartEmail := "/////@gmail.com"
	forbiddenCharactersInSecondPartEmail := "arman@////.///"

	// There is a length limit on email addresses.
	// That limit is a maximum of 64 characters (octets) in the "domain's first part" (right after the "@").
	// And a max of 64 chars right after the ".".
	tooLongInSecondPartEmail := "a@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com"
	tooLongInThirdPartEmail := "a@a.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	spaceCharacterInFirstPartEmail := "a a@a.a"
	spaceCharacterInSecondPartEmail := "a@a a.a"
	spaceCharacterInThirdPartEmail := "a@a.a a"
	missingDomainPartEmail := "email"
	missingDomainSecondPartEmail := "email@mail."
	missingUsernamePartEmail := "@mail.ru"
	goodEmail := "arman_alzhan@mail.ru"

	var tests = []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"ValidateEmail(specialCharactersInFirstPartEmail) must return true",
			specialCharactersInFirstPartEmail,
			true,
		},
		{
			"ValidateEmail(forbiddenCharactersInSecondPartEmail) must return false",
			forbiddenCharactersInSecondPartEmail,
			false,
		},
		{
			"ValidateEmail(tooLongInSecondPartEmail) must return false",
			tooLongInSecondPartEmail,
			false,
		},
		{
			"ValidateEmail(tooLongInThirdPartEmail) must return false",
			tooLongInThirdPartEmail,
			false,
		},
		{
			"ValidateEmail(spaceCharacterInFirstPartEmail) must return false",
			spaceCharacterInFirstPartEmail,
			false,
		},
		{
			"ValidateEmail(spaceCharacterInSecondPartEmail) must return false",
			spaceCharacterInSecondPartEmail,
			false,
		},
		{
			"ValidateEmail(spaceCharacterInThirdPartEmail) must return false",
			spaceCharacterInThirdPartEmail,
			false,
		},
		{
			"ValidateEmail(missingDomainPartEmail) must return false",
			missingDomainPartEmail,
			false,
		},
		{
			"ValidateEmail(missingUsernamePartEmail) must return false",
			missingUsernamePartEmail,
			false,
		},
		{
			"ValidateEmail(missingUsernamePartEmail) must return false",
			missingDomainSecondPartEmail,
			false,
		},
		{
			"ValidateEmail(goodEmail) must return false",
			goodEmail,
			true,
		},
	}

	for _, tst := range tests {
		v := validator.New()
		t.Run(tst.name, func(t *testing.T) {
			ValidateEmail(v, tst.input)
			result := v.Valid()
			if result != tst.expected {
				t.Errorf("Expected %v got %v", tst.expected, result)
			}
		})
	}
}

func TestValidatePasswordPlaintext(t *testing.T) {
	v := validator.New()

	emptyPassword := ""

	ValidatePasswordPlaintext(v, emptyPassword)
	expected := v.Valid()

	if expected != false {
		t.Errorf("ValidateProduct() returned unexpected value: got %v, expected %s", expected, "false")
	}
}

func TestTableDrivenValidatePasswordPlaintext(t *testing.T) {
	tooShortPassword := "hello"
	tooLongPassword := "ThisPasswordIsJustTooLongToBeAcceptedIntoThePasswordTableYouKnowThatIsForSureHowItIs"
	goodPassword := "goodPassword"
	var tests = []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"ValidatePasswordPlaintext(tooShortPassword) must return false",
			tooShortPassword,
			false,
		},
		{
			"ValidatePasswordPlaintext(tooLongPassword) must return false",
			tooLongPassword,
			false,
		},
		{
			"ValidatePasswordPlaintext(goodPassword) must return true",
			goodPassword,
			true,
		},
	}

	for _, tst := range tests {
		v := validator.New()
		t.Run(tst.name, func(t *testing.T) {
			ValidatePasswordPlaintext(v, tst.input)
			result := v.Valid()
			if result != tst.expected {
				t.Errorf("Expected %v got %v", tst.expected, result)
			}
		})
	}
}

func TestValidateUser(t *testing.T) {
	v := validator.New()

	noNameUser := &User{
		Name:        "",
		Surname:     "Surname",
		Username:    "username",
		DOB:         time.Now(),
		PhoneNumber: "+123456789",
		Address:     "Address",
		Email:       "arman_alzhan@mail.ru",
	}

	err := noNameUser.Password.Set("somePassword")
	if err != nil {
		t.Errorf("Password encryption returned an error. \n%v", err)
		return
	}

	ValidateUser(v, noNameUser)

	expected := v.Valid()

	if expected != false {
		t.Errorf("ValidateUser(noNameUser) returned unexpected value: got %v, expected %s", expected, "false")
	}
}
