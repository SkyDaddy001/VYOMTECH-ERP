package models

// Lead Status Constants - 30+ Detailed Statuses

// Initial Contact Phase - 11 statuses
const (
	LeadStatusFreshLead             = "Fresh Lead"
	LeadStatusReEngaged             = "Re Engaged"
	LeadStatusNotConnected          = "Not Connected"
	LeadStatusDead                  = "Dead"
	LeadStatusFollowUpCold          = "Follow Up - Cold"
	LeadStatusFollowUpWarm          = "Follow Up - Warm"
	LeadStatusFollowUpHot           = "Follow Up - Hot"
	LeadStatusLost                  = "Lost"
	LeadStatusUnqualifiedLocation   = "Unqualified (Location)"
	LeadStatusUnqualifiedBudget     = "Unqualified (Budget)"
	LeadStatusUnqualifiedClientProf = "Unqualified - Client Profile"
)

// Site Visit Phase - 12 statuses
const (
	LeadStatusSVScheduled        = "SV - Scheduled"
	LeadStatusSVDone             = "SV - Done"
	LeadStatusSVCold             = "SV - Cold"
	LeadStatusSVWarm             = "SV - Warm"
	LeadStatusSVRevisitScheduled = "SV - Revisit Scheduled"
	LeadStatusSVRevisitDone      = "SV - Revisit Done"
	LeadStatusSVHot              = "SV - Hot"
	LeadStatusSVLostNoResponse   = "SV - Lost (No Response)"
	LeadStatusSVLostBudget       = "SV - Lost (Budget)"
	LeadStatusSVLostPlanDropped  = "SV - Lost (Plan Dropped)"
	LeadStatusSVLostLocation     = "SV - Lost (Location)"
	LeadStatusSVLostAvailability = "SV - Lost (Availability)"
)

// Face-to-Face Phase - 10 statuses
const (
	LeadStatusF2FScheduled        = "F2F - Scheduled"
	LeadStatusF2FDone             = "F2F - Done"
	LeadStatusF2FFollowUp         = "F2F - Follow Up"
	LeadStatusF2FWarm             = "F2F - Warm"
	LeadStatusF2FHot              = "F2F - Hot"
	LeadStatusF2FLostNoResponse   = "F2F - Lost (No Response)"
	LeadStatusF2FLostBudget       = "F2F - Lost (Budget)"
	LeadStatusF2FLostPlanDropped  = "F2F - Lost (Plan Dropped)"
	LeadStatusF2FLostLocation     = "F2F - Lost (Location)"
	LeadStatusF2FLostAvailability = "F2F - Lost (Availability)"
)

// Booking Phase - 3 statuses
const (
	LeadStatusBookingInProgress   = "Booking-In-Progress"
	LeadStatusBookingProgressLost = "Booking Progress-Lost"
	LeadStatusBookingDone         = "Booking Done"
)

// Pipeline Stages - 7 stages
const (
	PipelineStageNew         = "New"
	PipelineStageConnected   = "Connected"
	PipelineStageInterested  = "Interested"
	PipelineStageAnalysis    = "Analysis"
	PipelineStageNegotiation = "Negotiation"
	PipelineStagePreBooking  = "Pre Booking"
	PipelineStageBooking     = "Booking"
)

// StatusToPipelineMap maps status to pipeline stage
var StatusToPipelineMap = map[string]string{
	// Initial Contact Phase
	LeadStatusFreshLead:             PipelineStageNew,
	LeadStatusReEngaged:             PipelineStageConnected,
	LeadStatusNotConnected:          PipelineStageNew,
	LeadStatusDead:                  PipelineStageNew,
	LeadStatusFollowUpCold:          PipelineStageConnected,
	LeadStatusFollowUpWarm:          PipelineStageConnected,
	LeadStatusFollowUpHot:           PipelineStageInterested,
	LeadStatusLost:                  PipelineStageConnected,
	LeadStatusUnqualifiedLocation:   PipelineStageConnected,
	LeadStatusUnqualifiedBudget:     PipelineStageConnected,
	LeadStatusUnqualifiedClientProf: PipelineStageConnected,

	// Site Visit Phase
	LeadStatusSVScheduled:        PipelineStageInterested,
	LeadStatusSVDone:             PipelineStageAnalysis,
	LeadStatusSVCold:             PipelineStageAnalysis,
	LeadStatusSVWarm:             PipelineStageAnalysis,
	LeadStatusSVRevisitScheduled: PipelineStageAnalysis,
	LeadStatusSVRevisitDone:      PipelineStageAnalysis,
	LeadStatusSVHot:              PipelineStageNegotiation,
	LeadStatusSVLostNoResponse:   PipelineStageAnalysis,
	LeadStatusSVLostBudget:       PipelineStageAnalysis,
	LeadStatusSVLostPlanDropped:  PipelineStageAnalysis,
	LeadStatusSVLostLocation:     PipelineStageAnalysis,
	LeadStatusSVLostAvailability: PipelineStageAnalysis,

	// Face-to-Face Phase
	LeadStatusF2FScheduled:        PipelineStageInterested,
	LeadStatusF2FDone:             PipelineStageAnalysis,
	LeadStatusF2FFollowUp:         PipelineStageAnalysis,
	LeadStatusF2FWarm:             PipelineStageAnalysis,
	LeadStatusF2FHot:              PipelineStageNegotiation,
	LeadStatusF2FLostNoResponse:   PipelineStageAnalysis,
	LeadStatusF2FLostBudget:       PipelineStageAnalysis,
	LeadStatusF2FLostPlanDropped:  PipelineStageAnalysis,
	LeadStatusF2FLostLocation:     PipelineStageAnalysis,
	LeadStatusF2FLostAvailability: PipelineStageAnalysis,

	// Booking Phase
	LeadStatusBookingInProgress:   PipelineStagePreBooking,
	LeadStatusBookingProgressLost: PipelineStagePreBooking,
	LeadStatusBookingDone:         PipelineStageBooking,
}

// AllLeadStatuses returns a list of all valid lead statuses
func AllLeadStatuses() []string {
	return []string{
		// Initial Contact
		LeadStatusFreshLead,
		LeadStatusReEngaged,
		LeadStatusNotConnected,
		LeadStatusDead,
		LeadStatusFollowUpCold,
		LeadStatusFollowUpWarm,
		LeadStatusFollowUpHot,
		LeadStatusLost,
		LeadStatusUnqualifiedLocation,
		LeadStatusUnqualifiedBudget,
		LeadStatusUnqualifiedClientProf,
		// Site Visit
		LeadStatusSVScheduled,
		LeadStatusSVDone,
		LeadStatusSVCold,
		LeadStatusSVWarm,
		LeadStatusSVRevisitScheduled,
		LeadStatusSVRevisitDone,
		LeadStatusSVHot,
		LeadStatusSVLostNoResponse,
		LeadStatusSVLostBudget,
		LeadStatusSVLostPlanDropped,
		LeadStatusSVLostLocation,
		LeadStatusSVLostAvailability,
		// Face-to-Face
		LeadStatusF2FScheduled,
		LeadStatusF2FDone,
		LeadStatusF2FFollowUp,
		LeadStatusF2FWarm,
		LeadStatusF2FHot,
		LeadStatusF2FLostNoResponse,
		LeadStatusF2FLostBudget,
		LeadStatusF2FLostPlanDropped,
		LeadStatusF2FLostLocation,
		LeadStatusF2FLostAvailability,
		// Booking
		LeadStatusBookingInProgress,
		LeadStatusBookingProgressLost,
		LeadStatusBookingDone,
	}
}

// IsValidLeadStatus checks if a status is valid
func IsValidLeadStatus(status string) bool {
	for _, s := range AllLeadStatuses() {
		if s == status {
			return true
		}
	}
	return false
}

// GetPipelineStage returns the pipeline stage for a given status
func GetPipelineStage(status string) string {
	if stage, ok := StatusToPipelineMap[status]; ok {
		return stage
	}
	return PipelineStageNew // Default to New
}
