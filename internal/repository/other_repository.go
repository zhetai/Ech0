package repository

import (
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/models"
)

func GetHeatMap(startDate, endDate string) ([]models.Heapmap, error) {
	var results []models.Heapmap

	// 查询数据
	// 执行查询
	err := database.DB.Table("messages").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", startDate, endDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
