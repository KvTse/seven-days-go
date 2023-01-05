package schema

import (
	"go-orm/dialect"
	"go/ast"
	"reflect"
)

// Field represents a column of table
type Field struct {
	Name string // 字段名
	Type string // 类型
	Tag  string // 约束
}

// Schema represents a table of database
type Schema struct {
	Model      interface{} // 对象
	Name       string      // 表名
	Fields     []*Field    // 属性
	FieldNames []string
	fieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i)
		if !f.Anonymous && ast.IsExported(f.Name) {
			field := &Field{Name: f.Name, Type: d.DataTypeOf(reflect.Indirect(reflect.New(f.Type)))}
			if value, ok := f.Tag.Lookup("geeorm"); ok {
				field.Tag = value
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, f.Name)
			schema.fieldMap[f.Name] = field
		}
	}

	return schema
}
