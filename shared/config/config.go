package config

import "time"

const (
	AccessTokenExpireDelay = time.Hour * 24 * 365
	RefreshTokenExpireDelay = time.Hour * 24 * 2
	SessionExpireTime = time.Hour * 24 * 30
)