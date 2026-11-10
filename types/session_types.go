package types

import (
	"time"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/config"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
)

type Session struct {
	Id        string    `bson:"_id,omitempty" json:"id,omitempty"`
	UserId    string    `bson:"userId" json:"userId"`
	TokenHash string    `bson:"tokenHash" json:"-"`
	ExpiresAt time.Time `bson:"expiresAt" json:"expiresAt"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

func NewSession(params *CreateSessionParams) *Session {
	now := time.Now()
	return &Session{
		UserId: params.UserId,
		TokenHash: params.TokenHash,
		ExpiresAt: now.Add(config.SessionExpireTime),
		CreatedAt: now,
	}
}

type CreateSessionParams struct {
	UserId    string `json:"userId"`
	TokenHash string `json:"-"`
}

func (params *CreateSessionParams) Validate() error {
	if len(params.UserId) == 0 {
		return localErrors.NewErrorBadRequest("userId is required")
	}
	if len(params.TokenHash) == 0 {
		return localErrors.NewErrorBadRequest("tokenHash is required")
	}
	return nil
}
