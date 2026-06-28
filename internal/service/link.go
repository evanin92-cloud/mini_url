package service

import (
	"errors"
	"time"

	"mini_url/internal/models"
	"mini_url/internal/repository"
	"mini_url/internal/utils"
)

type LinkService struct {
	linkRepo *repository.LinkRepository
}

func NewLinkService(linkRepo *repository.LinkRepository) *LinkService {
	return &LinkService{linkRepo: linkRepo}
}

func (s *LinkService) CreateLink(originalURL string, userID int, customID string) (*models.Link, error) {
	if !utils.ValidateURL(originalURL) {
		return nil, errors.New("invalid URL format")
	}

	var shortID string
	if customID != "" {
		shortID = customID
	} else {
		id, err := utils.GenerateRandomID(7)
		if err != nil {
			return nil, err
		}
		shortID = id
	}

	existing, _ := s.linkRepo.FindByShortID(shortID)
	if existing != nil {
		return nil, errors.New("short ID already exists")
	}

	link := &models.Link{
		OriginalURL: utils.NormalizeURL(originalURL),
		ShortID:     shortID,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ClickCount:  0,
		Status:      "active",
	}

	if err := s.linkRepo.Create(link); err != nil {
		return nil, err
	}

	return link, nil
}

func (s *LinkService) CreateBatch(urls []string, userID int) ([]*models.Link, error) {
	var links []*models.Link
	for _, url := range urls {
		link, err := s.CreateLink(url, userID, "")
		if err != nil {
			continue
		}
		links = append(links, link)
	}
	return links, nil
}

func (s *LinkService) GetLinksByUser(userID int) ([]*models.Link, error) {
	return s.linkRepo.FindByUserID(userID)
}

func (s *LinkService) DeleteLink(linkID, userID int) error {
	link, err := s.linkRepo.FindByID(linkID)
	if err != nil {
		return err
	}
	if link.UserID != userID {
		return errors.New("not authorized to delete this link")
	}
	return s.linkRepo.Delete(linkID)
}

func (s *LinkService) Redirect(shortID string) (string, error) {
	link, err := s.linkRepo.FindByShortID(shortID)
	if err != nil {
		return "", err
	}

	if link.Status != "active" {
		return "", errors.New("link is not active")
	}

	link.ClickCount++
	link.LastAccessedAt = time.Now()
	s.linkRepo.Update(link)

	return link.OriginalURL, nil
}