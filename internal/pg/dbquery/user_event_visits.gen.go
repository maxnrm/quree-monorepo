// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dbquery

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"quree/internal/pg/dbmodels"
)

func newUserEventVisit(db *gorm.DB, opts ...gen.DOOption) userEventVisit {
	_userEventVisit := userEventVisit{}

	_userEventVisit.userEventVisitDo.UseDB(db, opts...)
	_userEventVisit.userEventVisitDo.UseModel(&dbmodels.UserEventVisit{})

	tableName := _userEventVisit.userEventVisitDo.TableName()
	_userEventVisit.ALL = field.NewAsterisk(tableName)
	_userEventVisit.ID = field.NewString(tableName, "id")
	_userEventVisit.DateCreated = field.NewTime(tableName, "date_created")
	_userEventVisit.QuizID = field.NewString(tableName, "quiz_id")
	_userEventVisit.EventType = field.NewString(tableName, "event_type")
	_userEventVisit.UserID = field.NewString(tableName, "user_id")
	_userEventVisit.AdminID = field.NewString(tableName, "admin_id")

	_userEventVisit.fillFieldMap()

	return _userEventVisit
}

type userEventVisit struct {
	userEventVisitDo

	ALL         field.Asterisk
	ID          field.String
	DateCreated field.Time
	QuizID      field.String
	EventType   field.String
	UserID      field.String
	AdminID     field.String

	fieldMap map[string]field.Expr
}

func (u userEventVisit) Table(newTableName string) *userEventVisit {
	u.userEventVisitDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userEventVisit) As(alias string) *userEventVisit {
	u.userEventVisitDo.DO = *(u.userEventVisitDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userEventVisit) updateTableName(table string) *userEventVisit {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewString(table, "id")
	u.DateCreated = field.NewTime(table, "date_created")
	u.QuizID = field.NewString(table, "quiz_id")
	u.EventType = field.NewString(table, "event_type")
	u.UserID = field.NewString(table, "user_id")
	u.AdminID = field.NewString(table, "admin_id")

	u.fillFieldMap()

	return u
}

func (u *userEventVisit) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userEventVisit) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 6)
	u.fieldMap["id"] = u.ID
	u.fieldMap["date_created"] = u.DateCreated
	u.fieldMap["quiz_id"] = u.QuizID
	u.fieldMap["event_type"] = u.EventType
	u.fieldMap["user_id"] = u.UserID
	u.fieldMap["admin_id"] = u.AdminID
}

func (u userEventVisit) clone(db *gorm.DB) userEventVisit {
	u.userEventVisitDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u userEventVisit) replaceDB(db *gorm.DB) userEventVisit {
	u.userEventVisitDo.ReplaceDB(db)
	return u
}

type userEventVisitDo struct{ gen.DO }

type IUserEventVisitDo interface {
	gen.SubQuery
	Debug() IUserEventVisitDo
	WithContext(ctx context.Context) IUserEventVisitDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IUserEventVisitDo
	WriteDB() IUserEventVisitDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IUserEventVisitDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IUserEventVisitDo
	Not(conds ...gen.Condition) IUserEventVisitDo
	Or(conds ...gen.Condition) IUserEventVisitDo
	Select(conds ...field.Expr) IUserEventVisitDo
	Where(conds ...gen.Condition) IUserEventVisitDo
	Order(conds ...field.Expr) IUserEventVisitDo
	Distinct(cols ...field.Expr) IUserEventVisitDo
	Omit(cols ...field.Expr) IUserEventVisitDo
	Join(table schema.Tabler, on ...field.Expr) IUserEventVisitDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IUserEventVisitDo
	RightJoin(table schema.Tabler, on ...field.Expr) IUserEventVisitDo
	Group(cols ...field.Expr) IUserEventVisitDo
	Having(conds ...gen.Condition) IUserEventVisitDo
	Limit(limit int) IUserEventVisitDo
	Offset(offset int) IUserEventVisitDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IUserEventVisitDo
	Unscoped() IUserEventVisitDo
	Create(values ...*dbmodels.UserEventVisit) error
	CreateInBatches(values []*dbmodels.UserEventVisit, batchSize int) error
	Save(values ...*dbmodels.UserEventVisit) error
	First() (*dbmodels.UserEventVisit, error)
	Take() (*dbmodels.UserEventVisit, error)
	Last() (*dbmodels.UserEventVisit, error)
	Find() ([]*dbmodels.UserEventVisit, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*dbmodels.UserEventVisit, err error)
	FindInBatches(result *[]*dbmodels.UserEventVisit, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*dbmodels.UserEventVisit) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IUserEventVisitDo
	Assign(attrs ...field.AssignExpr) IUserEventVisitDo
	Joins(fields ...field.RelationField) IUserEventVisitDo
	Preload(fields ...field.RelationField) IUserEventVisitDo
	FirstOrInit() (*dbmodels.UserEventVisit, error)
	FirstOrCreate() (*dbmodels.UserEventVisit, error)
	FindByPage(offset int, limit int) (result []*dbmodels.UserEventVisit, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IUserEventVisitDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (u userEventVisitDo) Debug() IUserEventVisitDo {
	return u.withDO(u.DO.Debug())
}

func (u userEventVisitDo) WithContext(ctx context.Context) IUserEventVisitDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userEventVisitDo) ReadDB() IUserEventVisitDo {
	return u.Clauses(dbresolver.Read)
}

func (u userEventVisitDo) WriteDB() IUserEventVisitDo {
	return u.Clauses(dbresolver.Write)
}

func (u userEventVisitDo) Session(config *gorm.Session) IUserEventVisitDo {
	return u.withDO(u.DO.Session(config))
}

func (u userEventVisitDo) Clauses(conds ...clause.Expression) IUserEventVisitDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userEventVisitDo) Returning(value interface{}, columns ...string) IUserEventVisitDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userEventVisitDo) Not(conds ...gen.Condition) IUserEventVisitDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userEventVisitDo) Or(conds ...gen.Condition) IUserEventVisitDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userEventVisitDo) Select(conds ...field.Expr) IUserEventVisitDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userEventVisitDo) Where(conds ...gen.Condition) IUserEventVisitDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userEventVisitDo) Order(conds ...field.Expr) IUserEventVisitDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userEventVisitDo) Distinct(cols ...field.Expr) IUserEventVisitDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userEventVisitDo) Omit(cols ...field.Expr) IUserEventVisitDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userEventVisitDo) Join(table schema.Tabler, on ...field.Expr) IUserEventVisitDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userEventVisitDo) LeftJoin(table schema.Tabler, on ...field.Expr) IUserEventVisitDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userEventVisitDo) RightJoin(table schema.Tabler, on ...field.Expr) IUserEventVisitDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userEventVisitDo) Group(cols ...field.Expr) IUserEventVisitDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userEventVisitDo) Having(conds ...gen.Condition) IUserEventVisitDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userEventVisitDo) Limit(limit int) IUserEventVisitDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userEventVisitDo) Offset(offset int) IUserEventVisitDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userEventVisitDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IUserEventVisitDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userEventVisitDo) Unscoped() IUserEventVisitDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userEventVisitDo) Create(values ...*dbmodels.UserEventVisit) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userEventVisitDo) CreateInBatches(values []*dbmodels.UserEventVisit, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userEventVisitDo) Save(values ...*dbmodels.UserEventVisit) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userEventVisitDo) First() (*dbmodels.UserEventVisit, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*dbmodels.UserEventVisit), nil
	}
}

func (u userEventVisitDo) Take() (*dbmodels.UserEventVisit, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*dbmodels.UserEventVisit), nil
	}
}

func (u userEventVisitDo) Last() (*dbmodels.UserEventVisit, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*dbmodels.UserEventVisit), nil
	}
}

func (u userEventVisitDo) Find() ([]*dbmodels.UserEventVisit, error) {
	result, err := u.DO.Find()
	return result.([]*dbmodels.UserEventVisit), err
}

func (u userEventVisitDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*dbmodels.UserEventVisit, err error) {
	buf := make([]*dbmodels.UserEventVisit, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userEventVisitDo) FindInBatches(result *[]*dbmodels.UserEventVisit, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userEventVisitDo) Attrs(attrs ...field.AssignExpr) IUserEventVisitDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userEventVisitDo) Assign(attrs ...field.AssignExpr) IUserEventVisitDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userEventVisitDo) Joins(fields ...field.RelationField) IUserEventVisitDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userEventVisitDo) Preload(fields ...field.RelationField) IUserEventVisitDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userEventVisitDo) FirstOrInit() (*dbmodels.UserEventVisit, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*dbmodels.UserEventVisit), nil
	}
}

func (u userEventVisitDo) FirstOrCreate() (*dbmodels.UserEventVisit, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*dbmodels.UserEventVisit), nil
	}
}

func (u userEventVisitDo) FindByPage(offset int, limit int) (result []*dbmodels.UserEventVisit, count int64, err error) {
	result, err = u.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = u.Offset(-1).Limit(-1).Count()
	return
}

func (u userEventVisitDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userEventVisitDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userEventVisitDo) Delete(models ...*dbmodels.UserEventVisit) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userEventVisitDo) withDO(do gen.Dao) *userEventVisitDo {
	u.DO = *do.(*gen.DO)
	return u
}
