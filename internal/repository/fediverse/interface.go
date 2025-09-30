package repository

import (
	model "github.com/lin-snow/ech0/internal/model/fediverse"
)

type FediverseRepositoryInterface interface {
	GetFollowers(userID uint) ([]model.Follower, error)
}
