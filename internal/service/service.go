package service

import (
	"context"
	"fmt"
	"log"
	"service/internal/model"
	"service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CalculateCost(ctx context.Context, userID uuid.UUID, serviceName, startDate, endDate string) (int, error) {
	reqStart, err := time.Parse("01-2006", startDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start_date format: %w", err)
	}

	reqEnd, err := time.Parse("01-2006", endDate)
	if err != nil {
		return 0, fmt.Errorf("invalid end_date format: %w", err)
	}

	if reqStart.After(reqEnd) {
		return 0, fmt.Errorf("start_date cannot be after end_date")
	}

	subs, err := s.repo.GetForCostCalculation(ctx, userID, serviceName)
	if err != nil {
		return 0, err
	}

	totalCost := 0

	for _, sub := range subs {
		subStart, err := time.Parse("01-2006", sub.StartDate)
		if err != nil {
			continue // Игнорируем битые записи, чтобы не уронить весь расчет
		}

		subEnd := time.Now().AddDate(100, 0, 0) // Бесконечная подписка
		if sub.EndDate != nil {
			parsedEnd, err := time.Parse("01-2006", *sub.EndDate)
			if err == nil {
				subEnd = parsedEnd
			}
		}

		current := reqStart
		iterations := 0
		// Защита от бесконечного цикла (максимум 100 лет)
		for !current.After(reqEnd) && iterations < 1200 {
			if !current.Before(subStart) && !current.After(subEnd) {
				totalCost += sub.Price
			}
			current = current.AddDate(0, 1, 0)
			iterations++
		}
	}

	return totalCost, nil
}

func (s *Service) Create(ctx context.Context, sub model.Subscription) (int, error) {
    
    id, err := s.repo.Create(ctx, sub)
    if err != nil {
        
        log.Printf("Error creating subscription for service %s: %v", sub.ServiceName, err)
        return 0, err
    }

    log.Printf("Successfully created subscription with ID: %d", id)
    return id, nil
}