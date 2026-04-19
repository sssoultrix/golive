package domain

import "errors"

var ErrInvalidLogin = errors.New("invalid login")

var ErrInvalidCredentials = errors.New("invalid credentials")

var ErrInvalidRefreshToken = errors.New("invalid refresh token")

var ErrInvalidAccessToken = errors.New("invalid access token")
