package repository

import "time"

type RefreshToken struct {
	Token     string    `json:"token"`
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	IsRevoked bool      `json:"is_revoked"`
	CreatedAt time.Time `json:"created_at"`
	DeviceID  string    `json:"device_id,omitempty"`
}

type TokensPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type WebToken struct {
	AccessToken string `json:"access_token"`
}

type AuthRepository interface {
	StoreRefreshToken(rt RefreshToken) error
	GetRefreshToken(token string) (RefreshToken, error)
	InvalidateRefreshToken(token string) error
}
