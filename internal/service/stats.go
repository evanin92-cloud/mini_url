package service

import (
	"mini_url/internal/models"
	"mini_url/internal/repository"
)

type StatsService struct {
	statsRepo *repository.StatsRepository
}

func NewStatsService(statsRepo *repository.StatsRepository) *StatsService {
	return &StatsService{statsRepo: statsRepo}
}

func (s *StatsService) GetStats(shortID string) ([]*models.Stats, error) {
	return s.statsRepo.FindByShortID(shortID)
}