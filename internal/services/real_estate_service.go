package services

import "database/sql"

// RealEstateService provides real estate management functionality
type RealEstateService struct {
	DB *sql.DB
}

// NewRealEstateService creates a new real estate service instance
func NewRealEstateService(db *sql.DB) *RealEstateService {
	return &RealEstateService{
		DB: db,
	}
}
