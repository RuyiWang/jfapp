package mongo

import (
	"jfapp/pool"

	mgo "gopkg.in/mgo.v2"
)

type MongoSrc struct {
	*mgo.Session
}

var (
	MGO_CONN_CAP = 10
	session      *mgo.Session
	err          error
	MgoPool      = pool.ClassicPool(
		MGO_CONN_CAP,
		MGO_CONN_CAP/5,
		func() (pool.Src, error) {
			// if err != nil || session.Ping() != nil {
			// 	session, err = newSession()
			// }
			return &MongoSrc{session.Clone()}, err
		},
		60e9) //60e9 : 60秒
)

// 调用资源池中的资源
func Call(fn func(pool.Src) error) error {
	return MgoPool.Call(fn)
}

// 判断资源是否可用
func (self *MongoSrc) Usable() bool {
	if self.Session == nil || self.Session.Ping() != nil {
		return false
	}
	return true
}

// 使用后的重置方法
func (*MongoSrc) Reset() {}

// 被资源池删除前的自毁方法
func (self *MongoSrc) Close() {
	if self.Session == nil {
		return
	}
	self.Session.Close()
}
