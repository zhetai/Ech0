package service

import (
	repository "github.com/lin-snow/ech0/internal/repository/fediverse"
)

type FediverseService struct {
	fediverseRepository repository.FediverseRepositoryInterface
}

func NewFediverseService(fediverseRepo repository.FediverseRepositoryInterface) FediverseServiceInterface {
	return &FediverseService{
		fediverseRepository: fediverseRepo,
	}
}