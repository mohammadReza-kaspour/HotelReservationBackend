package types

import (
	"regexp"
	"time"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types/shared"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
)

type User struct {
	Id                string            `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string            `bson:"firstName" json:"firstName"`
	LastName          string            `bson:"lastName" json:"lastName"`
	Email             string            `bson:"email" json:"email"`
	EncryptedPassword string            `bson:"encryptedPassword" json:"-"`
	Role              []shared.RoleName `bson:"role" json:"role"`
	CreatedAt         time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time         `bson:"updatedAt" json:"updatedAt"`
}

func NewUser(params *CreateUserParams) (*User, error) {
	encryptPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		FirstName:         derefString(params.FirstName),
		LastName:          derefString(params.LastName),
		Email:             params.Email,
		EncryptedPassword: string(encryptPassword),
		Role:              []shared.RoleName{shared.RoleCustomer},
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

func (u *User) VerifyPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)); err != nil {
		return localErrors.NewErrorUnAuthorized("invalid password")
	}
	return nil
}

type CreateUserParams struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
}

func (params *CreateUserParams) Validate() error {
	validStringRegexp := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	validEmailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	validPasswordRegexp := regexp.MustCompile(`^[a-zA-Z0-9_.\-]{8,}$`)

	if params.FirstName != nil && !validStringRegexp.MatchString(*params.FirstName) {
		return localErrors.NewErrorBadRequest("user firstName is invalid")
	}
	if params.LastName != nil && !validStringRegexp.MatchString(*params.LastName) {
		return localErrors.NewErrorBadRequest("user lastName is invalid")
	}
	if !validEmailRegexp.MatchString(params.Email) {
		return localErrors.NewErrorBadRequest("user email is invalid")
	}
	if !validPasswordRegexp.MatchString(params.Password) {
		return localErrors.NewErrorBadRequest("user password is invalid")
	}

	return nil
}

type UpdateUserParams struct {
	Id        string  `bson:"_id" json:"id"`
	FirstName *string `bson:"firstName" json:"firstName"`
	LastName  *string `bson:"lastName" json:"lastName"`
	Email     *string `bson:"email" json:"email"`
}

func (params *UpdateUserParams) Validate() error {
	validStringRegexp := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	validEmailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if len(params.Id) == 0 {
		return localErrors.NewErrorBadRequest("user Id is required")
	}
	if params.FirstName != nil && !validStringRegexp.MatchString(*params.FirstName) {
		return localErrors.NewErrorBadRequest("user firstName is invalid")
	}
	if params.LastName != nil && !validStringRegexp.MatchString(*params.LastName) {
		return localErrors.NewErrorBadRequest("user lastName is invalid")
	}
	if params.Email != nil && !validEmailRegexp.MatchString(*params.Email) {
		return localErrors.NewErrorBadRequest("user email is invalid")
	}

	return nil
}
func (params *UpdateUserParams) ToUpdateMap() map[string]any {
	updateMap := map[string]any{}

	updateMap["id"] = params.Id
	derefNullable(updateMap, "firstName", params.FirstName)
	derefNullable(updateMap, "lastName", params.LastName)
	derefNullable(updateMap, "email", params.Email)

	if len(updateMap) == 0 {
		return nil
	}
	return updateMap
}

type LoginUserParams struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
}
func (params *LoginUserParams) Validate() error {
	validEmailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	validPasswordRegexp := regexp.MustCompile(`^[a-zA-Z0-9_.\-]{8,}$`)

	if !validEmailRegexp.MatchString(params.Email) {
		return localErrors.NewErrorBadRequest("user email is invalid")
	}
	if !validPasswordRegexp.MatchString(params.Password) {
		return localErrors.NewErrorBadRequest("user password is invalid")
	}

	return nil
}

type RefreshTokenParams struct {
	RefreshToken  string  `json:"refreshToken"`
}
func (params *RefreshTokenParams) Validate() error {
	if len(params.RefreshToken) == 0 {
		return localErrors.NewErrorBadRequest("refreshToken is required")
	}

	return nil
}