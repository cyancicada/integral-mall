package model

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"

	"integral-mall/common/baseerror"
)

type (
	//id           int(11)       (NULL)              NO      PRI     (NULL)             auto_increment  select,insert,update,references
	//good_name    varchar(200)  utf8mb4_general_ci  YES             (NULL)                             select,insert,update,references
	//price        int(11)       (NULL)              YES             (NULL)                             select,insert,update,references
	//intro        text          utf8mb4_general_ci  YES             (NULL)                             select,insert,update,references
	//image        varchar(255)  utf8mb4_general_ci  YES             (NULL)                             select,insert,update,references
	//store        int(11)       (NULL)              YES             (NULL)                             select,insert,update,references
	//create_time  timestamp     (NULL)              NO              CURRENT_TIMESTAMP                  select,insert,update,references
	Goods struct {
		Id         int64
		GoodName   string    `xorm:"varchar(200) notnull 'good_name'"`
		Price      int       `xorm:"int 'price'"`
		Intro      string    `xorm:"Text 'intro'"`
		Image      string    `xorm:"varchar(200) 'image'"`
		Store      int64     `xorm:"int   'store'"`
		CreateTime time.Time `xorm:"DateTime 'create_time'"`
	}

	GoodsModel struct {
		mysql      *xorm.Engine
		redisCache *redis.Client
		table      string
	}
)

const (
	goodsDefaultPageSize int = 10
)

var (
	ErrNotFound = baseerror.NewBaseError("没有找到相关的商品")
)

func NewGoodsModel(
	mysql *xorm.Engine,
	redisCache *redis.Client,
	table string,
) *GoodsModel {

	return &GoodsModel{mysql: mysql, redisCache: redisCache, table: table}
}

func (m *GoodsModel) Insert(u *Goods) (int64, error) {
	return m.mysql.Table(m.table).Insert(u)
}

func (m *GoodsModel) PageList(name string, page int) ([]*Goods, int64, error) {
	if page > 0 {
		page = page - 1
	}
	goodsList := []*Goods(nil)
	count, err := m.mysql.Table(m.table).Where("name LIKE ?", "%"+name+"%").Count(goodsList)
	if err != nil {
		return nil, 0, err
	}
	page = page * goodsDefaultPageSize
	if err := m.mysql.Table(m.table).Where("name LIKE ? LIMIT ?,10", "%"+name+"%", page).Find(&goodsList); err != nil {
		return nil, 0, err
	}
	return goodsList, count, nil

}

func (m *GoodsModel) FindById(id int64) (*Goods, error) {
	goods := new(Goods)
	if _, err := m.mysql.Table(m.table).Where("id = ?", id).Get(goods); err != nil {
		return nil, err
	}
	if goods.Id <= 0 {
		return nil, ErrNotFound
	}
	return goods, nil
}

func (m *GoodsModel) TransactionChangeStore(id, num int64, userId int, opts ...func(userId int) error) error {
	_, err := m.mysql.Transaction(func(session *xorm.Session) (i interface{}, e error) {
		query := "UPDATE " + m.table + " SET store=store-? WHERE id=?"
		if _, err := m.mysql.Exec(query, num, id); err != nil {
			return nil, err
		}
		for _, opt := range opts {
			if err := opt(userId); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	return err
}
