package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	// DataTypeOf 将go语言类型转换成对应数据库数据类型
	DataTypeOf(typ reflect.Value) string
	// TableExistSQL 检测某个表是否存在SQL语句
	TableExistSQL(tableName string) (string, []interface{})
}

// RegisterDialect 注册方法
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// 获取方言
func getDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
