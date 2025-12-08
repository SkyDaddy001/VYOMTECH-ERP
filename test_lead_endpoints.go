package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Test Lead API Endpoints
func testLeadAPI() {
	baseURL := "http://localhost:8080/api/v1"
	client := &http.Client{Timeout: 10 * time.Second}

	// Test data
	tenantID := "test-tenant-123"
	headers := map[string]string{
		"X-Tenant-ID":  tenantID,
		"Content-Type": "application/json",
	}

	fmt.Println("\n========================================")
	fmt.Println("LEAD API ENDPOINT TESTS")
	fmt.Println("========================================\n")

	// Test 1: Create a lead
	fmt.Println("[TEST 1] Create Lead - Fresh Lead Status")
	createLeadReq := map[string]interface{}{
		"name":   "John Doe",
		"email":  "john@example.com",
		"phone":  "+1234567890",
		"status": "Fresh Lead",
		"source": "manual",
		"notes":  "Potential customer",
	}

	reqBody, _ := json.Marshal(createLeadReq)
	req, _ := http.NewRequest("POST", baseURL+"/leads", bytes.NewBuffer(reqBody))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✓ Status: %d\n", resp.StatusCode)
		var lead map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&lead)
		leadID := int64(lead["id"].(float64))
		fmt.Printf("✓ Lead created with ID: %v\n\n", leadID)

		// Test 2: Update Lead Status to "Follow Up - Warm"
		fmt.Println("[TEST 2] Update Lead Status - Follow Up Warm")
		updateStatusReq := map[string]interface{}{
			"status": "Follow Up - Warm",
			"notes":  "Warm lead, schedule follow-up",
		}
		reqBody, _ := json.Marshal(updateStatusReq)
		req, _ := http.NewRequest("PUT", fmt.Sprintf(baseURL+"/leads/%d/status?id=%d", leadID, leadID), bytes.NewBuffer(reqBody))
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, _ := client.Do(req)
		fmt.Printf("✓ Status: %d\n", resp.StatusCode)
		var statusResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&statusResp)
		fmt.Printf("✓ Status Updated: %v\n", statusResp["status"])
		fmt.Printf("✓ Pipeline Stage: %v\n\n", statusResp["pipeline_stage"])

		// Test 3: Update to Site Visit Status
		fmt.Println("[TEST 3] Update Lead Status - Site Visit Done")
		updateStatusReq = map[string]interface{}{
			"status": "SV - Done",
			"notes":  "Site visit completed successfully",
		}
		reqBody, _ = json.Marshal(updateStatusReq)
		req, _ = http.NewRequest("PUT", fmt.Sprintf(baseURL+"/leads/%d/status?id=%d", leadID, leadID), bytes.NewBuffer(reqBody))
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, _ = client.Do(req)
		fmt.Printf("✓ Status: %d\n", resp.StatusCode)
		json.NewDecoder(resp.Body).Decode(&statusResp)
		fmt.Printf("✓ Status Updated: %v\n", statusResp["status"])
		fmt.Printf("✓ Pipeline Stage: %v\n\n", statusResp["pipeline_stage"])

		// Test 4: Update to F2F Meeting Status
		fmt.Println("[TEST 4] Update Lead Status - F2F Done")
		updateStatusReq = map[string]interface{}{
			"status": "F2F - Done",
			"notes":  "Face-to-face meeting completed",
		}
		reqBody, _ = json.Marshal(updateStatusReq)
		req, _ = http.NewRequest("PUT", fmt.Sprintf(baseURL+"/leads/%d/status?id=%d", leadID, leadID), bytes.NewBuffer(reqBody))
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, _ = client.Do(req)
		fmt.Printf("✓ Status: %d\n", resp.StatusCode)
		json.NewDecoder(resp.Body).Decode(&statusResp)
		fmt.Printf("✓ Status Updated: %v\n", statusResp["status"])
		fmt.Printf("✓ Pipeline Stage: %v\n\n", statusResp["pipeline_stage"])

		// Test 5: Update to Booking Status
		fmt.Println("[TEST 5] Update Lead Status - Booking In Progress")
		updateStatusReq = map[string]interface{}{
			"status": "Booking-In-Progress",
			"notes":  "Booking process initiated",
		}
		reqBody, _ = json.Marshal(updateStatusReq)
		req, _ = http.NewRequest("PUT", fmt.Sprintf(baseURL+"/leads/%d/status?id=%d", leadID, leadID), bytes.NewBuffer(reqBody))
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, _ = client.Do(req)
		fmt.Printf("✓ Status: %d\n", resp.StatusCode)
		json.NewDecoder(resp.Body).Decode(&statusResp)
		fmt.Printf("✓ Status Updated: %v\n", statusResp["status"])
		fmt.Printf("✓ Pipeline Stage: %v\n\n", statusResp["pipeline_stage"])

		// Test 6: Get Leads by Pipeline Stage
		fmt.Println("[TEST 6] Get Leads by Pipeline Stage - Booking")
		req, _ = http.NewRequest("GET", baseURL+"/leads/pipeline/Booking?stage=Booking", nil)
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, _ = client.Do(req)
		fmt.Printf("✓ Status: %d\n", resp.StatusCode)
		var leads []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&leads)
		fmt.Printf("✓ Leads in Booking stage: %d\n\n", len(leads))

		// Test 7: Get Leads by Status
		fmt.Println("[TEST 7] Get Leads by Specific Status - Booking In Progress")
		req, _ = http.NewRequest("GET", baseURL+"/leads/status/Booking-In-Progress?status=Booking-In-Progress", nil)
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, _ = client.Do(req)
		fmt.Printf("✓ Status: %d\n", resp.StatusCode)
		json.NewDecoder(resp.Body).Decode(&leads)
		fmt.Printf("✓ Leads with Booking-In-Progress status: %d\n\n", len(leads))

		// Test 8: Get All Leads
		fmt.Println("[TEST 8] Get All Leads")
		req, _ = http.NewRequest("GET", baseURL+"/leads", nil)
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, _ = client.Do(req)
		fmt.Printf("✓ Status: %d\n", resp.StatusCode)
		json.NewDecoder(resp.Body).Decode(&leads)
		fmt.Printf("✓ Total leads: %d\n\n", len(leads))

		// Test 9: Get Lead by ID
		fmt.Println(fmt.Sprintf("[TEST 9] Get Lead by ID - %d", leadID))
		req, _ = http.NewRequest("GET", fmt.Sprintf(baseURL+"/leads/%d?id=%d", leadID, leadID), nil)
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, _ = client.Do(req)
		fmt.Printf("✓ Status: %d\n", resp.StatusCode)
		var singleLead map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&singleLead)
		fmt.Printf("✓ Lead Name: %v\n", singleLead["name"])
		fmt.Printf("✓ Lead Status: %v\n", singleLead["status"])
		fmt.Printf("✓ Pipeline Stage: %v\n\n", singleLead["pipeline_stage"])

		// Test 10: Test Invalid Status
		fmt.Println("[TEST 10] Test Invalid Status (Should fail)")
		updateStatusReq = map[string]interface{}{
			"status": "InvalidStatus",
			"notes":  "This should fail",
		}
		reqBody, _ = json.Marshal(updateStatusReq)
		req, _ = http.NewRequest("PUT", fmt.Sprintf(baseURL+"/leads/%d/status?id=%d", leadID, leadID), bytes.NewBuffer(reqBody))
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, _ = client.Do(req)
		fmt.Printf("✓ Status: %d (expected 400)\n", resp.StatusCode)
		if resp.StatusCode == 400 {
			fmt.Println("✓ Invalid status properly rejected\n")
		}
	}

	fmt.Println("========================================")
	fmt.Println("ALL TESTS COMPLETED")
	fmt.Println("========================================\n")
}

// Test Status Constants
func testStatusConstants() {
	fmt.Println("\n========================================")
	fmt.Println("LEAD STATUS CONSTANTS VALIDATION")
	fmt.Println("========================================\n")

	statuses := map[string][]string{
		"Initial Contact": {
			"Fresh Lead", "Re Engaged", "Not Connected", "Dead",
			"Follow Up - Cold", "Follow Up - Warm", "Follow Up - Hot",
			"Lost", "Unqualified (Location)", "Unqualified (Budget)",
			"Unqualified - Client Profile",
		},
		"Site Visit": {
			"SV - Scheduled", "SV - Done", "SV - Cold", "SV - Warm",
			"SV - Revisit Scheduled", "SV - Revisit Done", "SV - Hot",
			"SV - Lost (No Response)", "SV - Lost (Budget)",
			"SV - Lost (Plan Dropped)", "SV - Lost (Location)",
			"SV - Lost (Availability)",
		},
		"Face-to-Face": {
			"F2F - Scheduled", "F2F - Done", "F2F - Follow Up", "F2F - Warm",
			"F2F - Hot", "F2F - Lost (No Response)", "F2F - Lost (Budget)",
			"F2F - Lost (Plan Dropped)", "F2F - Lost (Location)",
			"F2F - Lost (Availability)",
		},
		"Booking": {
			"Booking-In-Progress", "Booking Progress-Lost", "Booking Done",
		},
	}

	totalStatuses := 0
	for phase, statusList := range statuses {
		fmt.Printf("%s Phase: %d statuses\n", phase, len(statusList))
		totalStatuses += len(statusList)
	}

	fmt.Printf("\nTotal Statuses: %d\n\n", totalStatuses)

	pipelineStages := []string{
		"New", "Connected", "Interested", "Analysis",
		"Negotiation", "Pre Booking", "Booking",
	}

	fmt.Printf("Pipeline Stages: %d\n", len(pipelineStages))
	for i, stage := range pipelineStages {
		fmt.Printf("  %d. %s\n", i+1, stage)
	}

	fmt.Println("\n========================================")
	fmt.Println("VALIDATION COMPLETED")
	fmt.Println("========================================\n")
}
