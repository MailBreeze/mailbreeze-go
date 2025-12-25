package mailbreeze

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// AutomationsResource provides access to automation operations.
type AutomationsResource struct {
	client      *HTTPClient
	Enrollments *EnrollmentsResource
}

// Enroll enrolls a contact in an automation.
func (r *AutomationsResource) Enroll(ctx context.Context, params *EnrollParams) (*Enrollment, error) {
	var enrollment Enrollment
	if err := r.client.Post(ctx, "/automations/enroll", params, &enrollment); err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// EnrollmentsResource provides access to enrollment operations.
type EnrollmentsResource struct {
	client *HTTPClient
}

// List lists automation enrollments.
func (r *EnrollmentsResource) List(ctx context.Context, params *ListEnrollmentsParams) (*EnrollmentList, error) {
	query := url.Values{}

	if params != nil {
		if params.AutomationID != "" {
			query.Set("automation_id", params.AutomationID)
		}
		if params.Status != "" {
			query.Set("status", string(params.Status))
		}
		if params.Page > 0 {
			query.Set("page", strconv.Itoa(params.Page))
		}
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
	}

	var result EnrollmentList
	if err := r.client.Get(ctx, "/automations/enrollments", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Cancel cancels an enrollment.
func (r *EnrollmentsResource) Cancel(ctx context.Context, enrollmentID string) (*CancelEnrollmentResult, error) {
	var result CancelEnrollmentResult
	if err := r.client.Post(ctx, fmt.Sprintf("/automations/enrollments/%s/cancel", enrollmentID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
