package handler

import "errors"

var ErrorAuthorization = errors.New("no authorization header provided")

var ErrorBearer = errors.New("could not find bearer token in authorization header")
