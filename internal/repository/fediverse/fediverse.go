package repository

import (
	model "github.com/lin-snow/ech0/internal/model/fediverse"
	"gorm.io/gorm"
)

type FediverseRepository struct {
	db *gorm.DB
}

func NewFediverseRepository(db *gorm.DB) FediverseRepositoryInterface {
	return &FediverseRepository{
		db: db,
	}
}

func (r *FediverseRepository) GetFollowers(userID uint) ([]model.Follower, error) {
	var followers []model.Follower
	if err := r.db.Where("user_id = ?", userID).Find(&followers).Error; err != nil {
		return nil, err
	}
	return followers, nil
}

func (r *FediverseRepository) GetFollowing(userID uint) ([]model.Follow, error) {
	var following []model.Follow
	if err := r.db.Where("user_id = ?", userID).Find(&following).Error; err != nil {
		return nil, err
	}
	return following, nil
}