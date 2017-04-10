package orm

import (
	"github.com/go-xorm/xorm"
)

var OrmEngine *xorm.Engine

func NewSession() *xorm.Session {
	return OrmEngine.NewSession()
}

func CloseSession(sess *xorm.Session, err error) {
	if sess == nil {
		return
	}
	if err != nil {
		sess.Rollback()
	} else {
		sess.Commit()
	}
	sess.Close()
}

type BaseModel struct{}

type OrmModel interface {
	TableName() string
}

func (_ *BaseModel) ListByModel(sess *xorm.Session, modelSlicePtr interface{}, limit, offset int, orderBy string, orderType OrderType, condiBean OrmModel) (err error) {
	sess.Table(condiBean.TableName())
	if orderType == OrderTypeAsc {
		sess.Asc(orderBy)
	} else if orderType == OrderTypeDesc {
		sess.Desc(orderBy)
	}
	err = sess.Limit(limit, offset).Find(modelSlicePtr, condiBean)
	return
}

func (_ *BaseModel) ListByModelDescCreatedTime(sess *xorm.Session, modelSlicePtr interface{}, limit, offset int, condiBean OrmModel) (err error) {
	sess.Table(condiBean.TableName()).Desc("created")
	err = sess.Limit(limit, offset).Find(modelSlicePtr, condiBean)
	return
}

func (_ *BaseModel) ListByModelDescUpdatedTime(sess *xorm.Session, modelSlicePtr interface{}, limit, offset int, condiBean OrmModel) (err error) {
	sess.Table(condiBean.TableName()).Desc("updated")
	err = sess.Limit(limit, offset).Find(modelSlicePtr, condiBean)
	return
}

func (_ *BaseModel) CountByModel(sess *xorm.Session, condiBean OrmModel) (c int64, err error) {
	c, err = sess.Table(condiBean.TableName()).Count(condiBean)
	return
}
