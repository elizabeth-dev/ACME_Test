package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserDomainRepository(t *testing.T) {
	t.Parallel()

	for name, testGroup := range map[string]map[string]func(t *testing.T){
		"not found error": {
			"should return the correct error string": testNotFoundError,
		},
	} {
		testGroup := testGroup
		t.Run(
			name, func(t *testing.T) {
				for name, test := range testGroup {
					test := test
					t.Run(
						name, func(t *testing.T) {
							test(t)
						},
					)
				}
			},
		)
	}
}

func testNotFoundError(t *testing.T) {
	err := NotFoundError{Id: "1"}
	assert.Equal(t, "User with id 1 not found", err.Error())
}
