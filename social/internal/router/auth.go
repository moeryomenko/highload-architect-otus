package router

import (
	"encoding/json"
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
func (s *Auth) Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(BearerAuthScopes) == nil {
				next.ServeHTTP(w, r)
				return
			}

			token := extractHeaderToken(r)
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			var payload tokenPayload
			_, err := jwt.Verify([]byte(token), s.tokenizer, &payload)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			r.Header.Add(userIDHeaderKey, payload.UserID.String())
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
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
		UserID: userID,
	}
}

type tokenPayload struct {
	jwt.Payload
	UserID uuid.UUID `json:"user_id"`
}

// extractHeaderToken extracts jwt token from header.
func extractHeaderToken(r *http.Request) string {
	value := r.Header.Get(tokenHeader)
	return strings.TrimLeft(value, "Bearer ")
}
