package api

import (
	"context"
)

func (s *Handler) IsAuthenticated(ctx context.Context, token string) (bool, error) {
	_, err := parseVerifyToken(s.tokenSecret, token)
	if nil != err {
		return false, err
	}

	return true, nil
}
