package models

import "time"

type Tenant struct {
	Name      string    `json:"name"`
	Priority  int       `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
}
