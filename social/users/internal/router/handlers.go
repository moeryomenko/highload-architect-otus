package router

import (
	"encoding/json"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/moeryomenko/highload-architect-otus/social/internal/domain"
	"github.com/moeryomenko/highload-architect-otus/social/internal/repository"
	"github.com/moeryomenko/highload-architect-otus/social/internal/services"
)

type Service struct {
	logger *zap.Logger
	auth   *Auth
	login  *services.Login
	users  *repository.Users
}

// PostLogin logins user in service and generate access token.
func (s *Service) PostLogin(w http.ResponseWriter, r *http.Request) {
	var login PostLoginJSONRequestBody
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		ErrResponse(w, http.StatusBadRequest, err)
		return
	}
	l, err := s.login.Check(r.Context(), &domain.Login{Nickname: *login.Nickname, Password: *login.Password})
	switch err {
	case nil:
		s.auth.Sign(w, l.UserID)
	case services.ErrInvalidNickname, services.ErrInvalidPassword:
		ErrResponse(w, http.StatusBadRequest, err)
	default:
		ErrResponse(w, http.StatusInternalServerError, err)
	}
}

// PostSignup registrates new users and return assecc_token.
func (s *Service) PostSignup(w http.ResponseWriter, r *http.Request) {
	var login PostSignupJSONRequestBody
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	// TODO: move creation new user with id to application service.
	l := &domain.Login{
		UserID:   uuid.NewV4(),
		Nickname: *login.Nickname,
		Password: *login.Password,
	}
	err := s.login.EncryptSave(r.Context(), l)
	if err != nil {
		ErrResponse(w, http.StatusInternalServerError, err)
		return
	}
	s.auth.Sign(w, l.UserID)
}

// PutProfile saves user personal data.
func (s *Service) PutProfile(w http.ResponseWriter, r *http.Request) {
	var profileReq PutProfileJSONRequestBody
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&profileReq); err != nil {
		ErrResponse(w, http.StatusBadRequest, err)
		return
	}
	userID := ExtractUserID(r)

	profile := &domain.User{
		ID: userID,
		Info: &domain.Profile{
			FirstName: *profileReq.FirstName,
			LastName:  *profileReq.LastName,
			Age:       *profileReq.Age,
			Gender:    domain.Gender(*profileReq.Gender),
			Interests: *profileReq.Interests,
			City:      *profileReq.City,
		},
	}

	if err := s.users.Save(r.Context(), profile); err != nil {
		ErrResponse(w, http.StatusInternalServerError, err)
		return
	}
	resp, _ := json.Marshal(profileReq)
	w.Write(resp)
}

// GetProfiles returns paginated list of users.
func (s *Service) GetProfiles(w http.ResponseWriter, r *http.Request, params GetProfilesParams) {
	var pageOpts []repository.PageOption

	if params.PageSize != nil {
		pageOpts = append(pageOpts, repository.WithPageSize(*params.PageSize))
	}

	if params.PageToken != nil {
		token, err := repository.DecodeToken(*params.PageToken)
		if err != nil {
			ErrResponse(w, http.StatusBadRequest, err)
			return
		}
		pageOpts = append(pageOpts, repository.WithPageAt(token))
	}

	if params.Search != nil {
		if params.Search.FirstName != nil {
			pageOpts = append(pageOpts, repository.WithSearchByFirstName(*params.Search.FirstName))
		}
		if params.Search.LastName != nil {
			pageOpts = append(pageOpts, repository.WithSearchByLastName(*params.Search.LastName))
		}
	}

	page, nextToken, err := s.users.List(r.Context(), pageOpts...)
	if err != nil {
		s.logger.Error("could not get profiles list", zap.Error(err))
		ErrResponse(w, http.StatusInternalServerError, err)
		return
	}
	profiles := make([]Profile, 0, len(page))
	for _, profile := range page {
		profiles = append(profiles, MapProfileToResp(&profile))
	}

	resp, err := json.Marshal(Profiles{Profiles: &profiles, NextPageToken: &nextToken})
	if err != nil {
		ErrResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(resp)
}

func MapProfileToResp(user *domain.User) Profile {
	return Profile{
		FirstName: &user.Info.FirstName,
		LastName:  &user.Info.LastName,
		Age:       &user.Info.Age,
		Gender:    (*string)(&user.Info.Gender),
		Interests: &user.Info.Interests,
		City:      &user.Info.City,
	}
}

// ErrResponse writes error response with specified status.
func ErrResponse(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	errs, _ := json.Marshal(Error{Errors: &[]string{err.Error()}})
	w.Write(errs)
}
