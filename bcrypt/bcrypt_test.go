package bcrypt

import (
	"testing"
)

const BcryptCost = 12
const TestPassword = "testPassword"

func TestAll(t *testing.T) {
	hashed, err := Encrypt(TestPassword, BcryptCost)
	if err != nil {
		panic(err)
	}
	t.Log("Hashed: ", hashed)

	same, err := Compare(TestPassword, hashed)
	if !same {
		panic(err)
	}
	t.Log("(Original, hashed): ", TestPassword, ", ", hashed)
	return
}
