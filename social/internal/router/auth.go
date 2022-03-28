package router

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	jwt "github.com/gbrlsnchs/jwt/v3"
	uuid "github.com/satori/go.uuid"

	"github.com/moeryomenko/highload-architect-otus/social/internal/config"
)

const (
	tokenHeader     = "Authorization"
	userIDHeaderKey = "user_id"
)

// ExtractUserID is helper function for http handlers extract user id from request context.
func ExtractUserID(r *http.Request) uuid.UUID {
	id, _ := uuid.FromString(r.Header.Get(userIDHeaderKey))
	return id
}

// Auth is auth service.
type Auth struct {
	tokenizer  *jwt.HMACSHA
	expiration time.Duration
}

// NewAuth returns instance of Auth service.
func NewAuth(cfg *config.Config) *Auth {
	return &Auth{
		tokenizer: jwt.NewHS256([]byte(cfg.Session.Secret)),
	}
}

// Sign sets auth token to response.
func (s *Auth) Sign(w http.ResponseWriter, userID uuid.UUID) error {
	token, err := jwt.Sign(s.genToken(userID), s.tokenizer)
	if err != nil {
		return err
	}

	tokenStr := string(token)
	resp, err := json.Marshal(&Token{AccessToken: &tokenStr})
	if err != nil {
		return err
	}

	_, err = w.Write(resp)
	return err
}

// Auth is middleware that verify auth token and extract user id to request context.
func (s *Auth) Auth() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(BearerAuthScopes) == nil {
				next(w, r)
				return
			}
			payload, err := s.verify(r)
			if err != nil {
				ErrResponse(w, http.StatusUnauthorized, err)
				return
			}
			r.Header.Add(userIDHeaderKey, payload.UserID)
			next(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// verify verifies request and return access token.
func (s *Auth) verify(r *http.Request) (*tokenPayload, error) {
	token := extractHeaderToken(r)
	if token == "" {
		return nil, errors.New("dont have token")
	}
	payload := &tokenPayload{}
	_, err := jwt.Verify([]byte(token), s.tokenizer, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// genToken generates token payload.
func (s *Auth) genToken(userID uuid.UUID) *tokenPayload {
	now := time.Now()
	return &tokenPayload{
		Payload: jwt.Payload{
			Issuer:         "social",
			Audience:       jwt.Audience{},
			ExpirationTime: jwt.NumericDate(now.Add(s.expiration)),
			IssuedAt:       jwt.NumericDate(now),
		},
		UserID: userID.String(),
	}
}

type tokenPayload struct {
	jwt.Payload
	UserID string `json:"user_id"`
}

// extractHeaderToken extracts jwt token from header.
func extractHeaderToken(r *http.Request) string {
	value := r.Header.Get(tokenHeader)
	values := strings.Split(value, " ")
	if len(values) != 2 {
		return ""
	}
	return values[1]
}
