package models

import (
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        *primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string              `json:"name" bson:"name" validate:"required,min=2,max=50"`
	Email     string              `json:"email" bson:"email" validate:"required,email"`
	Password  string              `json:"password" bson:"password" validate:"required,min=8,max=100"`
	Provider  *string             `json:"provider" bson:"provider"`
	Role      *string             `json:"role" bson:"role"`
	Photo     *string             `json:"photo" bson:"photo"`
	Verified  *bool               `json:"verified" bson:"verified"`
	CreatedAt *time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time          `json:"updated_at" bson:"updated_at"`
}

type SignUpInput struct {
	Name            string `json:"name" bson:"name" validate:"required,min=2,max=50"`
	Email           string `json:"email" bson:"email" validate:"required,email"`
	Password        string `json:"password" bson:"password" validate:"required,min=8,max=100"`
	PasswordConfirm string `json:"password_confirm" bson:"password_confirm" validate:"required,min=8,max=100"`
	Photo           string `json:"photo" bson:"photo"`
}

type SignInInput struct {
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required,min=8,max=100"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty" `
	Role      string             `json:"role,omitempty" bson:"role,omitempty"`
	Photo     string             `json:"photo,omitempty" bson:"photo,omitempty"`
	Provider  string             `json:"provider,omitempty" bson:"provider,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func FilteredUserResponse(user *User) UserResponse {
	return UserResponse{
		ID:        *user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      *user.Photo,
		Provider:  *user.Provider,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, &ErrorResponse{
				Field: err.Field(),
				Tag:   err.Tag(),
				Value: err.Param(),
			})

		}
	}
	return errors
}
