package models

import "time"

type ReservationStatus string

const (
	StatusActive    ReservationStatus = "active"
	StatusExpired   ReservationStatus = "expired"
	StatusPreempted ReservationStatus = "preempted"
	StatusReleased  ReservationStatus = "released"
)

type Reservation struct {
	ID         string            `json:"id"`
	TenantName string            `json:"tenant_name"`
	PoolName   string            `json:"pool_name"`
	Resources  Resources         `json:"resources"`
	Status     ReservationStatus `json:"status"`
	TTL        int               `json:"ttl"`
	CreatedAt  time.Time         `json:"created_at"`
	ExpiresAt  *time.Time        `json:"expires_at,omitempty"`
	ReleasedAt *time.Time        `json:"released_at,omitempty"`
}

func (r *Reservation) IsExpired() bool {
	if r.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*r.ExpiresAt)
}

func (r *Reservation) IsActive() bool {
	return r.Status == StatusActive && !r.IsExpired()
}
