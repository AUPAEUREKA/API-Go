package model

import (
	"errors"
	"time"

	old "github.com/bearbin/go-age"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UUID        string    `json:"uuid"`
	AccessLevel int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	DateOfBirth time.Time `json:"birth_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u User) Valid(db *gorm.DB) []error {
	var errs []error
	if len(u.Password) == 0 {
		errs = append(errs, errors.New("No given password"))
	}
	if old.Age(u.DateOfBirth) < 18 {
		errs = append(errs, errors.New("You are not adult!"))
	}
	r := User{}
	_ = db.Where("email = ?", r.Email).First(&r)
	if r.Email == u.Email {
		errs = append(errs, errors.New("User with this email already exist"))
	}
	if len(errs) != 0 {
		return errs
	}
	return nil
}

// BeforeCreate : Gorm hook
func (u *User) BeforeCreate(scope *gorm.Scope) {
	id, _ := uuid.NewV4()
	u.UUID = id.String()
	scope.SetColumn("Password", hashPassword(u.Password))
	scope.SetColumn("UUID", u.UUID)
	return
}

// AfterFind : Remove password from the user to avoid security issues
func (u *User) AfterFind() (err error) {
	u.Password = ""
	return
}

// AfterCreate : Remove password from the user to avoid security issues
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	u.Password = ""
	return
}

// hashPassword : simple password hashing method
func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

// CheckPasswordHash : Compare password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
