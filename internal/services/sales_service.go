package services

import "database/sql"

type SalesService struct {
	DB *sql.DB
}

// NewSalesService creates a new sales service
func NewSalesService(db *sql.DB) *SalesService {
	return &SalesService{
		DB: db,
	}
}
