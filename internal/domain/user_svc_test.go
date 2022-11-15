package domain

import (
	"context"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	testcases := map[string]struct {
		u        User
		password string
		err      string
	}{
		"success": {
			u: User{
				Email:  "test@test.com",
				Name:   "test",
				Gender: GenderMale,
				Age:    18,
			},
			password: "123456",
		},
		"no email": {
			u: User{
				Name:   "test",
				Gender: GenderMale,
				Age:    18,
			},
			password: "123456",
			err:      "invalid email",
		},
		"bad email": {
			u: User{
				Email:  "test.com",
				Name:   "test",
				Gender: GenderMale,
				Age:    18,
			},
			password: "123456",
			err:      "invalid email",
		},
		"no name": {
			u: User{
				Email:  "test@test.com",
				Gender: GenderMale,
				Age:    18,
			},
			password: "123456",
			err:      "empty name",
		},
		"no gender": {
			u: User{
				Email: "test@test.com",
				Name:  "test",
				Age:   18,
			},
			password: "123456",
			err:      "unknown gender",
		},
		"unknown gender": {
			u: User{
				Email:  "test@test.com",
				Name:   "test",
				Gender: "test",
				Age:    18,
			},
			password: "123456",
			err:      "unknown gender",
		},
		"age lower than 18": {
			u: User{
				Email:  "test@test.com",
				Name:   "test",
				Gender: GenderMale,
				Age:    17,
			},
			password: "123456",
			err:      "age must be equal or above 18",
		},
		"short password": {
			u: User{
				Email:  "test@test.com",
				Name:   "test",
				Gender: GenderMale,
				Age:    18,
			},
			password: "12456",
			err:      "password must have at least 6 characters",
		},
		"no password": {
			u: User{
				Email:  "test@test.com",
				Name:   "test",
				Gender: GenderMale,
				Age:    18,
			},
			err: "password must have at least 6 characters",
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			userSvc := NewUserSvc(userRepoStub{})
			_, err := userSvc.Create(context.Background(), tc.u, tc.password)
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if tc.err == "" && errMsg != "" {
				t.Errorf("expected no error but received %q", errMsg)
			}
			if !strings.Contains(errMsg, tc.err) {
				t.Errorf("\nexpected %q\nreceived %q", tc.err, errMsg)
			}
		})
	}
}

func TestCreateRandom(t *testing.T) {
	userSvc := NewUserSvc(userRepoStub{654})
	u, password, err := userSvc.CreateRandom(context.Background())
	if err != nil {
		t.Fatalf("expected no error: %v", err)
	}

	if u.ID != 654 {
		t.Fatal("the id returned from the repo has not been assigned to the user")
	}

	if password == "" {
		t.Fatal("missing password")
	}

	if err := u.Validate(); err != nil {
		t.Fatalf("expected a valid user: %v", err)
	}
}

type userRepoStub struct {
	id uint32
}

func (r userRepoStub) Create(ctx context.Context, u User, password []byte) (uint32, error) {
	return r.id, nil
}

func (r userRepoStub) Get(ctx context.Context, userID uint32) (user User, err error) {
	return User{}, nil
}

func (r userRepoStub) ListPotentialMatches(ctx context.Context, user User) ([]User, error) {
	return nil, nil
}
