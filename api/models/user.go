package models

import (
	"time"

	"../core"
	"github.com/astaxie/beego/orm"
)

// User User
type User struct {
	BaseModel

	Nickname     string        `orm:"column(nickname);" json:"nickname,omitempty"`
	Congregation *Congregation `orm:"column(congregation_id);rel(fk);null" json:"congregation,omitempty"`
	Email        string        `orm:"column(email);" json:"email,omitempty"`
	Phone        string        `orm:"column(phone);" json:"phone,omitempty"`
	Name         string        `orm:"column(name);" json:"name,omitempty"`
	Password     string        `orm:"column(password);" json:"password,omitempty"`
	Role         string        `orm:"column(role);" json:"role,omitempty"`
	LastActivity *time.Time    `orm:"column(last_activity);type(timestamp)" json:"lastActivity,omitempty"`
}

// TableName TableName
func (t *User) TableName() string {
	return "users"
}

func init() {
	orm.RegisterModel(new(User))
}

// LoadCongregation loads congregation
func (t *User) LoadCongregation() {
	o := orm.NewOrm()
	o.LoadRelated(t, "Congregation", 0)
}

// RenewAuthToken - renews auth token
func (t *User) RenewAuthToken() (*AuthToken, error) {
	authToken := &AuthToken{}
	qs := core.GetModelQuerySeter(nil, new(AuthToken), true)
	err := qs.Filter("user_id__id", t.ID).One(authToken)
	if err != nil {
		authToken = NewAuthToken(t)
		_, err = core.InsertModel(nil, authToken)
		if err != nil {
			return nil, err
		}
	} else {
		authToken.LastLogin = time.Now()
		err = core.UpdateModel(nil, authToken, []string{"last_login"})
		if err != nil {
			return nil, err
		}
	}
	return authToken, nil
}

// GetUserByNickname returns user by nickname
func GetUserByNickname(nickname string) (*User, error) {
	user := new(User)
	err := core.GetModelQuerySeter(nil, user, false).Filter("nickname", nickname).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
