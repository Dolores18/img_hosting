package dao

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

// 利用泛型和反射优化查询
func FindByFields[T any](db *gorm.DB, model T, conditions map[string]interface{}, isOrQuery bool) ([]map[string]interface{}, error) {
	var results []T
	query := db.Model(&model)

	if len(conditions) > 0 {
		for field, value := range conditions {
			if isOrQuery {
				query = query.Or(fmt.Sprintf("%s = ?", field), value)
			} else {
				query = query.Where(fmt.Sprintf("%s = ?", field), value)
			}
		}
	}

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	var mapResults []map[string]interface{}
	for _, result := range results {
		mapResult := make(map[string]interface{})
		value := reflect.ValueOf(result)
		typ := reflect.TypeOf(result)
		for i := 0; i < value.NumField(); i++ {
			field := typ.Field(i)
			mapResult[field.Name] = value.Field(i).Interface()
		}
		mapResults = append(mapResults, mapResult)
	}

	return mapResults, nil
}
