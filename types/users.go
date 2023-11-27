package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 3
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() []string {
	errs := []string{}
	if len(params.FirstName) < minFirstNameLen {
		e := fmt.Sprintf("\"firstName\" %d karakteden asagi ola bilmez", minFirstNameLen)
		errs = append(errs, e)
	}
	if len(params.LastName) < minLastNameLen {
		e := fmt.Sprintf("\"lastName\" %d karakteden asagi ola bilmez", minLastNameLen)
		errs = append(errs, e)
	}
	if len(params.Password) < minPasswordLen {
		e := fmt.Sprintf("\"password\" %d karakteden asagi ola bilmez", minPasswordLen)
		errs = append(errs, e)
	}
	if !isEmailValid(params.Email) {
		e := fmt.Sprintf("\"email\" Adresi formata uygun deyil, %s", params.Email)
		errs = append(errs, e)
	}
	if len(errs)>0{
		return errs
	}
	return nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9._%+\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

type Users struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*Users, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Users{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
