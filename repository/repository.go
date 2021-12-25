package repository

import (
	"math"

	"github.com/sing3demons/go-todos/model"
	"gorm.io/gorm"
)

func paginate(value interface{}, pagination *model.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	ch := make(chan int64)
	go countRecords(ch, value, db)

	pagination.TotalRows = <-ch
	totalPages := int(math.Ceil(float64(pagination.TotalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.Limit).Order(pagination.GetSort())
	}
}

func countRecords(ch chan int64, value interface{}, db *gorm.DB) {
	var totalRows int64
	db.Model(value).Count(&totalRows)
	ch <- totalRows
}
