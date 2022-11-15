package domain

import (
	"math/rand"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	testcases := map[string]struct {
		u   User
		err string
	}{
		"success": {
			u: User{
				Email:  "test@test.com",
				Name:   "test",
				Gender: GenderMale,
				Age:    18,
			},
		},
		"no email": {
			u: User{
				Name:   "test",
				Gender: GenderMale,
				Age:    18,
			},
			err: "invalid email",
		},
		"bad email": {
			u: User{
				Email:  "test.com",
				Name:   "test",
				Gender: GenderMale,
				Age:    18,
			},
			err: "invalid email",
		},
		"no name": {
			u: User{
				Email:  "test@test.com",
				Gender: GenderMale,
				Age:    18,
			},
			err: "empty name",
		},
		"no gender": {
			u: User{
				Email: "test@test.com",
				Name:  "test",
				Age:   18,
			},
			err: "unknown gender",
		},
		"unknown gender": {
			u: User{
				Email:  "test@test.com",
				Name:   "test",
				Gender: "test",
				Age:    18,
			},
			err: "unknown gender",
		},
		"age lower than 18": {
			u: User{
				Email:  "test@test.com",
				Name:   "test",
				Gender: GenderMale,
				Age:    17,
			},
			err: "age must be equal or above 18",
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			userSvc := UserSvc{r: userRepoStub{}}
			_, err := userSvc.Create(tc.u)
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if !strings.Contains(errMsg, tc.err) {
				t.Errorf("test failed:\nexpected %q\nreceived %q", tc.err, errMsg)
			}
		})
	}
}

func TestCreateRandom(t *testing.T) {
	userSvc := UserSvc{r: userRepoStub{654}}
	u, err := userSvc.CreateRandom()
	if err != nil {
		t.Fatalf("expected no error: %v", err)
	}

	if u.ID != 654 {
		t.Fatal("the id returned from the repo has not been assigned to the user")
	}

	if err := u.Validate(); err != nil {
		t.Fatalf("expected a valid user: %v", err)
	}
}

type userRepoStub struct {
	id uint32
}

func (r userRepoStub) Create(u User) (id uint32, err error) {
	if id != 0 {
		return id, nil
	}

	return rand.Uint32(), nil
}
