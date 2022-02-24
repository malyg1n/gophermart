package request

// AuthRequest login and register request.
type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
