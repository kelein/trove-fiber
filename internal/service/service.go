package service

import (
	"github.com/kelein/trove-gin/pkg/jwt"
	"github.com/kelein/trove-gin/pkg/sid"

	"github.com/kelein/trove-fiber/internal/repository"
)

// Service stands for backend service layer
type Service struct {
	sid *sid.Sid
	jwt *jwt.JWT
	tm  repository.Transaction
}

// NewService creates a new Service instance.
func NewService(sid *sid.Sid, jwt *jwt.JWT, tm repository.Transaction) *Service {
	return &Service{
		sid: sid,
		jwt: jwt,
		tm:  tm,
	}
}
