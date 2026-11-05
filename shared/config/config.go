package config

import "time"

const (
	AccessTokenExpireDelay = time.Hour * 50_000
	RefreshTokenExpireDelay = time.Hour * 24
)