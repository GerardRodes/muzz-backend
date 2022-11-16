package domain

import (
	"context"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
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

			userSvc := NewService(userRepoStub{}, nil)
			_, err := userSvc.CreateUser(context.Background(), tc.u, tc.password)
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

func TestCreateRandomUser(t *testing.T) {
	userSvc := NewService(userRepoStub{654}, nil)
	u, password, err := userSvc.CreateRandomUser(context.Background())
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

func (r userRepoStub) GetUser(ctx context.Context, userID uint32) (User, error) {
	return User{}, nil
}

func (r userRepoStub) CreateUser(ctx context.Context, user User, passwordHash []byte) (uint32, error) {
	return r.id, nil
}

func (r userRepoStub) ListPotentialMatches(ctx context.Context, user User) ([]User, error) {
	return nil, nil
}

func (r userRepoStub) Swipe(ctx context.Context, userID, profileID uint32, preference bool) error {
	return nil
}

func (r userRepoStub) BothLiked(ctx context.Context, userID, profileID uint32) (bool, error) {
	return false, nil
}

func (r userRepoStub) CreateMatch(ctx context.Context, userID1, userID2 uint32) (uint64, error) {
	return 0, nil
}

func (r userRepoStub) GetUserIDAndPasswordByEmail(ctx context.Context, email string) (userID uint32, passHash []byte, err error) {
	return 0, nil, nil
}
