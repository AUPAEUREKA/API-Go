package model

import (
	"github.com/jinzhu/gorm"
)

// Proposal : the proposal struct definition
type User struct {
	gorm.Model
	UUID        string `json:"uuid"`
	AccessLevel int    `json:"access"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	DateOfBirth int    `json:"birth_date"`
}

//type JsonBirthDate time.Time

// implement Marshaler und Unmarshalere interface
/*func (j *JsonBirthDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = JsonBirthDate(t)
	return nil
}

func (j JsonBirthDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(j)
}

// Maybe a Format function for printing your date
func (j JsonBirthDate) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}

/*func (u User) Valid(db *gorm.DB) []error {
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
}*/

// BeforeCreate : Gorm hook
/*func (u *User) BeforeCreate(scope *gorm.Scope) {
	id, _ := uuid.NewV4()
	u.UUID = id.String()
	scope.SetColumn("Password", hashPassword(u.Password))
	scope.SetColumn("UUID", u.UUID)
	return
}

// hashPassword : simple password hashing method
/*func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

// CheckPasswordHash : Compare password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}*/
