package model

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type (
	User struct {
		Id         int64
		Name       string    `xorm:"varchar(20) notnull 'name'"`
		Mobile     string    `xorm:"varchar(25) notnull  unique 'mobile'"`
		Password   string    `xorm:"varchar(25) notnull  'password'"`
		CreateTime time.Time `xorm:"DateTime 'create_time'"`
	}

	UserModel struct {
		mysql      *xorm.Engine
		redisCache *redis.Client
		table      string
	}
)

const (
	AuthorizationExpire = 604800 * time.Second
)

func NewUserModel(
	mysql *xorm.Engine,
	redisCache *redis.Client,
	table string,
) *UserModel {

	return &UserModel{mysql: mysql, redisCache: redisCache, table: table}
}

func (m *UserModel) Insert(u *User) (int64, error) {
	return m.mysql.Insert(u)
}

func (m *UserModel) TransactionInsert(u *User, opts ...func(userId int64) error) (*User, error) {
	_, err := m.mysql.Transaction(func(session *xorm.Session) (i interface{}, e error) {
		if _, err := m.mysql.Insert(u); err != nil {
			return nil, err
		}
		for _, opt := range opts {
			if err := opt(u.Id); err != nil {
				return nil, err
			}
		}
		return u, nil
	})
	return u, err
}

func (m *UserModel) ExistByMobile(mobile string) (bool, error) {
	return m.mysql.Exist(&User{Mobile: mobile})
}

func (m *UserModel) FindByMobile(mobile string) (*User, error) {
	u := new(User)
	if _, err := m.mysql.Where("mobile = ?", mobile).Get(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (m *UserModel) FindById(id int64) (*User, error) {
	u := new(User)
	if _, err := m.mysql.Where("id = ?", id).Get(u); err != nil {
		return nil, err
	}
	return u, nil
}
