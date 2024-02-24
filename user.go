package musclemem

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// User represents a registered user
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

// userService implements methods for /users
type userService service

func (s *userService) Login() {}

func (s *userService) Register(ctx context.Context, u *User) (*User, *http.Response, error) {
	path := "/users"

	userJSON, err := json.Marshal(u)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.send(ctx, http.MethodPost, path, bytes.NewReader(userJSON))
	if err != nil {
		return nil, nil, err
	}

	respUser := new(User)
	if err := json.NewDecoder(resp.Body).Decode(&respUser); err != nil {
		return nil, nil, err
	}

	return respUser, resp, nil
}
