package types

import (
	"golang.org/x/crypto/bcrypt"
)

const cost = 12

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Password string `json:"password"`
	Confirm string `json:"confirm"`
	Validator
}

type UpdateUserParams struct {
	FirstName string `json:"firstName,omitempty"`
	LastName string `json:"lastName,omitempty"`
	Password string `json:"password,omitempty"`
	Confirm string `json:"confirm,omitempty"`
	Validator
}

type UpdatedUser struct{
	FirstName string `bson:"firstName,omitempty"`
	LastName string `bson:"lastName,omitempty"`
	Password string `bson:"password,omitempty"`
}

type User struct {
	ID string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName string `bson:"lastName" json:"lastName"`
	Email string `bson:"email" json:"email"`
	IsAdmin bool `bson:"isAdmin" json:"-"`
	Password string `bson:"password" json:"-"`
}

type LoginParams struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Validator
}

func NewUserParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), cost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName: params.FirstName,
		LastName: params.LastName,
		Email: params.Email,
		IsAdmin: false,
		Password: string(encpw),
	}, nil
}

func NewUpdatedUser(params UpdateUserParams) (*UpdatedUser, error) {
	
	u := &UpdatedUser{
		FirstName: params.FirstName,
		LastName: params.LastName,
	}

	if params.Password != "" {
		pw, err := bcrypt.GenerateFromPassword([]byte(params.Password), cost)
		if err != nil {
			return nil, err
		}
		u.Password = string(pw)
	}
	return u, nil
}

func (p *CreateUserParams) Validate() bool{ 
	p.CheckField(NotBlank(p.FirstName), "name", "This field cannot be blank")
	p.CheckField(NotBlank(p.LastName), "surname", "This field cannot be blank")
	p.CheckField(NotBlank(p.Email), "email", "This field cannot be blank")
	p.CheckField(Matches(p.Email, EmailRX), "email", "This field must be a valid email address")
	p.CheckField(NotBlank(p.Password), "password", "This field cannot be blank")
	p.CheckField(MinChars(p.Password, 8), "password", "This field must be at least 8 characters long")
	p.CheckField(NotBlank(p.Confirm), "confirm", "This field cannot be blank")
	p.CheckField(ConfirmPassword(p.Password, p.Confirm), "confirm", "Passwords do not match")

	return p.Valid()
}

func (p *UpdateUserParams) Validate() bool{
	if p.FirstName != ""{
		p.CheckField(NotBlank(p.FirstName), "name", "This field cannot be blank")
	}
	if p.LastName != ""{
		p.CheckField(NotBlank(p.LastName), "surname", "This field cannot be blank")
	}
	if p.Password != "" && p.Confirm != ""{
		p.CheckField(NotBlank(p.Password), "password", "This field cannot be blank")
		p.CheckField(MinChars(p.Password, 8), "password", "This field must be at least 8 characters long")
		p.CheckField(NotBlank(p.Confirm), "confirm", "This field cannot be blank")
		p.CheckField(ConfirmPassword(p.Password, p.Confirm), "confirm", "Passwords do not match")
	}else if (p.Password != "" && p.Confirm == "") || (p.Password == "" && p.Confirm != ""){
		p.Validator.AddFieldError("password", "in order to update password fill out both password and confirm")
	}
	return p.Valid()
}

func (p *LoginParams) Validate() bool {
	p.CheckField(NotBlank(p.Email), "email", "This field cannot be blank")
	p.CheckField(NotBlank(p.Password), "password", "This field cannot be blank")
	
	return p.Valid()
}