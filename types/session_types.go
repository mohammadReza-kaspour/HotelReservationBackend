package types

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/config"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
)

type Session struct {
	Id        string    `bson:"_id,omitempty" json:"id,omitempty"`
	UserId    string    `bson:"userId" json:"userId"`
	TokenHash string    `bson:"tokenHash" json:"-"`
	Revoked   bool      `bson:"revoked" json:"revoked"`
	ExpiresAt time.Time `bson:"expiresAt" json:"expiresAt"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

func NewSession(params *CreateSessionParams) *Session {
	now := time.Now()
	return &Session{
		UserId:    params.UserId,
		TokenHash: params.hashToken(),
		Revoked:   false,
		ExpiresAt: now.Add(config.SessionExpireTime),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type CreateSessionParams struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

func (params *CreateSessionParams) Validate() error {
	if len(params.UserId) == 0 {
		return localErrors.NewErrorBadRequest("userId is required")
	}
	if len(params.Token) == 0 {
		return localErrors.NewErrorBadRequest("tokenHash is required")
	}
	return nil
}

func (params *CreateSessionParams) hashToken() string {
	sum := sha256.Sum256([]byte(params.Token))
	return hex.EncodeToString(sum[:])
}

type UpdateSessionParams struct {
	Token   string `json:"token"`
	Revoked *bool  `json:"revoked"`
}

func (params *UpdateSessionParams) Validate() error {
	if len(params.Token) == 0 {
		return localErrors.NewErrorBadRequest("tokenHash is required")
	}
	return nil
}

func (params *UpdateSessionParams) ToUpdateMap() map[string]any {
	updateMap := map[string]any{}

	updateMap["tokenHash"] = params.hashToken()
	derefNullable(updateMap, "revoked", params.Revoked)

	if len(updateMap) == 0 {
		return nil
	}
	return updateMap
}

func (params *UpdateSessionParams) hashToken() string {
	sum := sha256.Sum256([]byte(params.Token))
	return hex.EncodeToString(sum[:])
}
