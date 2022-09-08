package domain

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	id        string
	name      string
	birthDate time.Time
	sex       Sex
	email     string
	password  string
	isActive  bool

	info    *UserInfo
	options *UserOptions
}

func (u *User) IsActive() bool {
	return u.isActive
}

func (u *User) Info() *UserInfo {
	return u.info
}

func (u *User) Options() *UserOptions {
	return u.options
}

func (u *User) SetInfo(info *UserInfo) {
	u.info = info
}

func (u *User) SetOptions(options *UserOptions) {
	u.options = options
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Sex() Sex {
	return u.sex
}

func (u *User) BirthDate() time.Time {
	return u.birthDate
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}

func (u *User) Id() string {
	return u.id
}

type UserInfo struct {
	about     string
	instagram string
	geo       *Geo
}

func (u *UserInfo) About() string {
	return u.about
}

func (u *UserInfo) Instagram() string {
	return u.instagram
}

func (u *UserInfo) Geo() *Geo {
	return u.geo
}

func NewUserInfo(about string, instagram string, geo *Geo) *UserInfo {
	return &UserInfo{about: about, instagram: instagram, geo: geo}
}

type UserOptions struct {
	sexPreference Sex
	distance      int
}

func (u *UserOptions) SexPreference() Sex {
	return u.sexPreference
}

func (u *UserOptions) Distance() int {
	return u.distance
}

func NewUserOptions(sexPreference Sex, distance int) *UserOptions {
	return &UserOptions{sexPreference: sexPreference, distance: distance}
}

func NewUser(name string, email string, password string, birthDate time.Time, sex Sex) (*User, error) {
	newPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	return &User{name: name, email: email, password: string(newPass), isActive: true, birthDate: birthDate, sex: sex}, nil
}

func NewUserFromDb(id string, name string, birthDate time.Time, email string, password string, isActive bool, sex Sex) *User {
	return &User{id: id, name: name, birthDate: birthDate, email: email, password: password, isActive: isActive, sex: sex}
}

func (u *User) ComparePassword(notDecodedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.password), []byte(notDecodedPassword))
}

func (u *User) ToJSON() interface{} {

	return struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Instagram string `json:"instagram"`
		About     string `json:"about"`
		Sex       Sex    `json:"sex"`
		SexPref   Sex    `json:"sex_preference"`
		Distance  int    `json:"distance"`
		Geo       struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"geo"`
		BirthDate string `json:"birth_date"`
		IsActive  bool   `json:"isActive"`
	}{
		u.id,
		u.name,
		u.email,
		u.info.instagram,
		u.info.about,
		u.sex,
		u.options.sexPreference,
		u.options.distance,
		struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		}{
			u.info.geo.latitude.String(),
			u.info.geo.longitude.String(),
		},
		u.birthDate.Format("02-01-2006"),
		u.isActive,
	}
}
