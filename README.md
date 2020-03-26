# integral-mall

# sql 如下 


```sql```

  CREATE DATABASE `goods` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;
CREATE DATABASE `integral` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;
CREATE DATABASE `order` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;
CREATE DATABASE `user` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

CREATE TABLE goods.goods
(
  id          int(11)       NOT NULL AUTO_INCREMENT primary key,
  good_name   varchar(255)  not null default '' comment '商品名',
  price       decimal(5, 2) not null default 0 comment '价格',
  intro       text comment '详情',
  image       varchar(255)  not null default '' comment '商品图片',
  store       int(11)       not null default 0 comment '库存',
  create_time timestamp     NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
  update_time timestamp     NOT NULL DEFAULT current_timestamp on update current_timestamp
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;


CREATE TABLE integral.integral
(
  id          int(11)   NOT NULL AUTO_INCREMENT primary key,
  user_id     int       not null default 0 comment '用户id',
  integral    int       not null default 0 comment '积分',
  create_time timestamp NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT current_timestamp on update current_timestamp
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;


CREATE TABLE order.order
(
  id          int(11)      NOT NULL AUTO_INCREMENT primary key,
  user_id     int          not null default 0 comment '用户id',
  goods_id    int          not null default 0 comment '商品id',
  good_name   varchar(255) NOT NULL default '' COMMENT '名称',
  mobile      varchar(15)  NOT NULL DEFAULT '' COMMENT '手机号',
  num         varchar(60)  NOT NULL DEFAULT '' COMMENT '数量',
  create_time timestamp    NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
  update_time timestamp    NOT NULL DEFAULT current_timestamp on update current_timestamp
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;

CREATE TABLE user.user
(
  id          int(11)     NOT NULL AUTO_INCREMENT primary key,
  name        varchar(11) NOT NULL default '' COMMENT '用户名称',
  mobile      varchar(15) NOT NULL DEFAULT '' COMMENT '手机号',
  password    varchar(60) NOT NULL DEFAULT '' COMMENT '密码',
  create_time timestamp   NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
  update_time timestamp   NOT NULL DEFAULT current_timestamp on update current_timestamp
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8


```sql```
