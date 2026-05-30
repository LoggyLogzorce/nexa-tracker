package auth

import "time"

type SessionResponse struct {
	ID         uint      `json:"id"`
	DeviceName string    `json:"device_name"`
	IPAddress  *string   `json:"ip_address"`
	LastActive time.Time `json:"last_active"`
	IsCurrent  bool      `json:"is_current"`
}
