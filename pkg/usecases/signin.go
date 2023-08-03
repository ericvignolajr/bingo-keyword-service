package usecases

import "github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"

type SigninRequest struct {
	Email string
}
type SigninResponse struct {
	Ok  bool
	Err error
}

type Signin struct {
	UserStore inmemory.UserStore
}

func (s *Signin) Exec(r SigninRequest) SigninResponse {
	_, err := s.UserStore.ReadByEmail(r.Email)
	if err != nil {
		return SigninResponse{
			false,
			err,
		}
	}

	return SigninResponse{
		true,
		nil,
	}

	// THIS MAY BE A GOOD PLACE TO SEPARATE USER CREATION AND USER RETRIEVAL
	// INTO SEPARATE USECASES SO THAT THEY CAN EVOLVE SEPARATELY
	// THE CONTROLLER WOULD USE THE CREATE USER USECASE IF THE READ SIGNIN
	// USECASE FAILS
	// user does not exist
	// 1. create a new user
	// 2. persist to user store
}
