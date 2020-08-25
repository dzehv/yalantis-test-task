package models

import (
	"log"
	"os"
	"strings"

	u "yal/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// JWT perms
type Token struct {
	UserId uint
	jwt.StandardClaims
}

type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
	About    string `json:"about"`
}

// validate user data
func (account *Account) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is too short"), false
	}

	// email should be uniq
	temp := new(Account)

	defer CloseDB()
	// check duplicates
	if err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error; err != nil &&
		err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry! "+err.Error()), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user!"), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (account *Account) Create() map[string]interface{} {
	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	defer CloseDB()
	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error")
	}

	// new JWT for new account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	// no password in response
	account.Password = ""

	resp := u.Message(true, "Account has been created")
	resp["account"] = account

	return resp
}

func Login(email, password string) map[string]interface{} {
	account := new(Account)

	defer CloseDB()
	if err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password),
		[]byte(password)); err != nil &&
		err == bcrypt.ErrMismatchedHashAndPassword {
		// creds incorrect
		return u.Message(false, "Invalid login credentials")
	}
	// no password in response
	account.Password = ""

	// create token
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	// save token to resp
	account.Token = tokenString

	resp := u.Message(true, "Logged In")
	resp["account"] = account

	return resp
}

func GetUser(id uint) map[string]interface{} {
	acc := new(Account)

	defer CloseDB()
	GetDB().Table("accounts").Where("id = ?", id).First(acc)

	// no user found
	if acc.Email == "" {
		return nil
	}
	acc.Password = ""

	resp := u.Message(true, "Got user")
	resp["account"] = acc

	return resp
}

func GetUsers() map[string]interface{} {
	accs := make([]Account, 0)

	defer CloseDB()
	if err := GetDB().Table("accounts").Find(&accs).Error; err != nil {
		log.Printf("Request error: %v", err)
		return nil
	}

	if len(accs) < 1 {
		log.Printf("No users found in DB")
		return nil
	}

	resp := u.Message(true, "Got user")
	resp["accounts"] = accs

	return resp
}

func (account *Account) Update() map[string]interface{} {
	// TODO: validate and allow to update only self account

	// update password hash
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	// save all fields
	defer CloseDB()
	GetDB().Save(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to update account, connection error")
	}

	// also delete password
	account.Password = ""

	resp := u.Message(true, "Account has been updated")
	resp["account"] = account

	return resp
}
