package services

import (
    "fixitpolytech/internal/models"
    "time"
)

type RequestService struct {
    requests []models.Request
}

func NewRequestService() *RequestService {
    return &RequestService{
        requests: make([]models.Request, 0),
    }
}

func (s *RequestService) CreateRequest(buildingID int, fieldID int, additionalText string, photos []string) models.Request {
    request := models.Request{
        BuildingID:    buildingID,
        FieldID:       fieldID,
        AdditionalText: additionalText,
        Photos:        photos,
        Status:        "not taken",
        Time:          time.Now(),
    }
    s.requests = append(s.requests, request)
    return request
}

func (s *RequestService) GetRequests() []models.Request {
    return s.requests
}