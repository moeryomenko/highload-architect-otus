package router

import (
	"net/http"
	"testing"
	"time"

	jwt "github.com/gbrlsnchs/jwt/v3"
	"github.com/stretchr/testify/assert"
)

func TestAuth_verify(t *testing.T) {
	tests := []struct {
		name string
		args *http.Request
		want string
	}{
		{
			name: "valid token",
			args: &http.Request{Header: http.Header{tokenHeader: []string{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzb2NpYWwiLCJleHAiOjE2NDg0NTQ0OTgsImlhdCI6MTY0ODQ1NDQ5OCwidXNlcl9pZCI6IjgzOWVhNzBlLTYxMzYtNDVhMy1hZjIxLWMwNDUxZDgzMDkzYyJ9.R2Ki9JigIPg76RJhwVVeC_3OSfARvTpdu4kKfBvwOXE"}}},
			want: "839ea70e-6136-45a3-af21-c0451d83093c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Auth{
				tokenizer:  jwt.NewHS256([]byte("secret")),
				expiration: 5 * time.Minute,
			}
			token, err := s.verify(tt.args)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, token.UserID)
		})
	}
}
