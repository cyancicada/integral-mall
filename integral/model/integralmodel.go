package model

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type (
	Integral struct {
		Id         int64
		UserId     int       `xorm:"int unique 'user_id'"`
		Integral   int       `xorm:"int  'integral'"`
		CreateTime time.Time `xorm:"DateTime 'create_time'"`
	}

	IntegralModel struct {
		mysql      *xorm.Engine
		redisCache *redis.Client
		table      string
	}
)

func NewIntegralModel(
	mysql *xorm.Engine,
	redisCache *redis.Client,
	table string,
) *IntegralModel {

	return &IntegralModel{mysql: mysql, redisCache: redisCache, table: table}
}

func (m *IntegralModel) Insert(data *Integral) (int64, error) {
	return m.mysql.Insert(data)
}

func (m *IntegralModel) ExistByUserId(userId int) (bool, error) {
	return m.mysql.Exist(&Integral{UserId: userId})
}

func (m *IntegralModel) FindByUserId(userId int) (*Integral, error) {
	res := new(Integral)
	if _, err := m.mysql.Where("userId = ?", userId).Get(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *IntegralModel) UpdateIntegralByUserId(userId, integral int) (*Integral, error) {
	query := "update `" + m.table + "` set `integral`=`integral`-? where user_id=?"
	if _, err := m.mysql.Exec(query, integral, userId); err != nil {
		return nil, err
	}
	return m.FindByUserId(userId)
}
func (m *IntegralModel) FindById(id int64) (*Integral, error) {
	res := new(Integral)
	if _, err := m.mysql.Where("id = ?", id).Get(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *IntegralModel) InsertIntegralSql(userId, integral int) string {
	return fmt.Sprintf("INSERT INTO "+m.table+" (`user_id`,`integral`) VALUES (%d,%d)", userId, integral)
}

func (m *IntegralModel) UpdateIntegralByUserIdSql(userId, integral int) string {
	query := "update `" + m.table + "` set `integral`=`integral`-? where user_id=?"
	return fmt.Sprintf(query, userId, integral)
}

func (m *IntegralModel) ExecSql(sql string) error {
	if _, err := m.mysql.Exec(sql); err != nil {
		return err
	}
	return nil
}
