package auth

import "net/http"

func GetBearerToken(headers http.Header) (string, error) {}

// chromium --ozone-platform-hint=wayland
