package models

import "time"

// ConstructionProject represents a construction project
type ConstructionProject struct {
	ID                     int64  `gorm:"primaryKey"`
	TenantID               string `gorm:"index"`
	ProjectName            string
	ProjectCode            string `gorm:"uniqueIndex"`
	Location               string
	Client                 string
	ContractValue          float64
	StartDate              time.Time
	ExpectedCompletion     time.Time
	CurrentProgressPercent int
	Status                 string // planning, active, suspended, completed, on_hold
	ProjectManager         string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// BillOfQuantities represents items in the bill of quantities
type BillOfQuantities struct {
	ID              int64  `gorm:"primaryKey"`
	TenantID        string `gorm:"index"`
	ProjectID       int64
	BOQNumber       string
	ItemDescription string
	Unit            string
	Quantity        float64
	UnitRate        float64
	TotalAmount     float64
	Category        string // civil, structural, electrical, plumbing, finishing, other
	Status          string // planned, in_progress, completed, on_hold
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ProgressTracking represents project progress records
type ProgressTracking struct {
	ID                int64  `gorm:"primaryKey"`
	TenantID          string `gorm:"index"`
	ProjectID         int64
	Date              time.Time
	ActivityDesc      string `gorm:"type:text"`
	QuantityCompleted float64
	Unit              string
	PercentComplete   int
	WorkforceDeployed int
	Notes             string `gorm:"type:text"`
	PhotoURL          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// QualityControl represents quality control inspections
type QualityControl struct {
	ID               int64  `gorm:"primaryKey"`
	TenantID         string `gorm:"index"`
	ProjectID        int64
	BOQItemID        int64
	InspectionDate   time.Time
	InspectorName    string
	QualityStatus    string // passed, failed, partial, pending
	Observations     string `gorm:"type:text"`
	CorrectiveAction string `gorm:"type:text"`
	FollowUpDate     time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ConstructionEquipment represents equipment used in construction
type ConstructionEquipment struct {
	ID             int64  `gorm:"primaryKey"`
	TenantID       string `gorm:"index"`
	ProjectID      int64
	EquipmentName  string
	EquipmentType  string
	SerialNumber   string
	Status         string // available, in_use, maintenance, retired
	DeploymentDate time.Time
	RetirementDate time.Time
	CostPerDay     float64
	Notes          string `gorm:"type:text"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TableName specifies the table name for ConstructionProject
func (ConstructionProject) TableName() string {
	return "construction_projects"
}

// TableName specifies the table name for BillOfQuantities
func (BillOfQuantities) TableName() string {
	return "bill_of_quantities"
}

// TableName specifies the table name for ProgressTracking
func (ProgressTracking) TableName() string {
	return "progress_tracking"
}

// TableName specifies the table name for QualityControl
func (QualityControl) TableName() string {
	return "quality_control"
}

// TableName specifies the table name for ConstructionEquipment
func (ConstructionEquipment) TableName() string {
	return "construction_equipment"
}
