package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmailAndName(t *testing.T) {
	tests := []struct {
		Name        string
		User        User
		ExpectError bool
	}{
		{
			Name:        "Valid Name and Email",
			User:        User{Name: "One1 yean", Email: "test@gmail.com"},
			ExpectError: false,
		},
		{
			Name:        "Missing Name",
			User:        User{Name: "", Email: "test@gmail.com"},
			ExpectError: false,
		},
		{
			Name:        "Missing Email",
			User:        User{Name: "One1 yean", Email: ""},
			ExpectError: false,
		},
		{
			Name:        "Missing Both Name and Email",
			User:        User{Name: "", Email: ""},
			ExpectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := test.User.ValidateEmailAndName()

			if test.ExpectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
