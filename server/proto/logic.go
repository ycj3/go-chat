package proto

type LoginRequest struct {
	UserID string
}

type LoginResponse struct {
	Code      int
	AuthToken string
}
