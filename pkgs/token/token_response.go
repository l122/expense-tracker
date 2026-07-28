package token

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID           string `json:"id"`
		Email        string `json:"email"`
		UserMetadata struct {
			AvatarURL     string `json:"avatar_url"`
			FullName      string `json:"full_name"`
			EmailVerified bool   `json:"email_verified"`
		} `json:"user_metadata"`
	} `json:"user"`
}
