package model

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type (
	Order struct {
		Id         string    `xorm:"varchar(100) notnull 'id'"`
		GoodName   string    `xorm:"varchar(20) notnull 'good_name'"`
		Mobile     string    `xorm:"varchar(25) notnull  unique 'mobile'"`
		Num        int64     `xorm:"int   'num'"`
		UserId     int64     `xorm:"int   'user_id'"`
		GoodsId    int64     `xorm:"int   'goods_id'"`
		CreateTime time.Time `xorm:"DateTime 'create_time'"`
	}

	OrderModel struct {
		mysql      *xorm.Engine
		redisCache *redis.Client
		table      string
	}
)

const (
	OrderDefaultPageSize int = 10
)

func NewOrderModel(
	mysql *xorm.Engine,
	redisCache *redis.Client,
	table string,
) *OrderModel {

	return &OrderModel{mysql: mysql, redisCache: redisCache, table: table}
}

func (m *OrderModel) Insert(o *Order) (int64, error) {
	return m.mysql.Insert(o)
}

func (m *OrderModel) FindById(id string) (*Order, error) {
	u := new(Order)
	if _, err := m.mysql.Where("id = ?", id).Get(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (m *OrderModel) PageFindByUserId(userId int64, page int) ([]*Order, int64, error) {
	orders := []*Order(nil)
	var count int64
	var err error
	if count, err = m.mysql.Where("user_id = ?", userId).Count(&orders); err != nil {
		return nil, 0, err
	}
	if page > 0 {
		page = page - 1
	}
	start := OrderDefaultPageSize * page
	if err = m.mysql.Where("user_id = ?", userId).Limit(OrderDefaultPageSize, start).Find(&orders); err != nil {
		return nil, 0, err
	}
	return orders, count, nil
}

func (m *OrderModel) ExecSql(sql string) error {
	if _, err := m.mysql.Exec(sql); err != nil {
		return err
	}
	return nil
}

func (m *OrderModel) BookingGoodsSql(orderId string, goodsId int64, goodsName, mobile string, userId, num int64) string {
	return fmt.Sprintf(
		"INSERT INTO "+m.table+" (id,goods_id,good_name,mobile,num,user_id) VALUES ('%s',%d,'%s','%s',%d,%d)",
		orderId, goodsId, goodsName, mobile, num, userId,
	)
}
