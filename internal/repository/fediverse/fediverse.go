package repository

import "gorm.io/gorm"

type FediverseRepository struct {
	db *gorm.DB
}

func NewFediverseRepository(db *gorm.DB) FediverseRepositoryInterface {
	return &FediverseRepository{
		db: db,
	}
}