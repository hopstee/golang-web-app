package service

import (
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"
)

type RequestService struct {
	requestRepo repository.RequestRepository
	logger      *slog.Logger
}

func NewRequestService(requestRepo repository.RequestRepository, logger *slog.Logger) *RequestService {
	return &RequestService{
		requestRepo: requestRepo,
		logger:      logger,
	}
}

func (s *RequestService) Get() ([]repository.Request, error) {
	s.logger.Info("get requests attempt")
	requests, err := s.requestRepo.Get()
	if err != nil {
		s.logger.Warn("get requests failed: requests not found", slog.Any("err", err))
		return nil, err
	}

	s.logger.Info("get requests successfull")
	return requests, nil
}

func (s *RequestService) GetByID(id int64) (repository.Request, error) {
	s.logger.Info("get request attempt")
	request, err := s.requestRepo.GetByID(id)
	if err != nil {
		s.logger.Warn("get request failed: request not found", slog.Any("err", err))
		return repository.Request{}, err
	}

	s.logger.Info("get requests successfull")
	return request, nil
}

func (s *RequestService) Create(request repository.Request) (int64, error) {
	s.logger.Info("create request attempt")
	id, err := s.requestRepo.Create(request)
	if err != nil {
		s.logger.Warn("create request failed", slog.Any("err", err))
		return 0, err
	}

	s.logger.Info("create request successfull")
	return id, nil
}

func (s *RequestService) Update(request repository.Request) error {
	s.logger.Info("update request attempt")
	err := s.requestRepo.Update(request)
	if err != nil {
		s.logger.Warn("update request failed: request not found", slog.Any("err", err))
		return err
	}

	s.logger.Info("update request successfull")
	return nil
}

func (s *RequestService) Delete(id int64) error {
	s.logger.Info("delete request attempt")
	err := s.requestRepo.Delete(id)
	if err != nil {
		s.logger.Warn("delete request failed: request not found", slog.Any("err", err))
		return err
	}

	s.logger.Info("delete request successfull")
	return nil
}
