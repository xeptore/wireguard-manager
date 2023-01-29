package api

import (
	"context"
)

func (h *Handler) IsAuthenticated(ctx context.Context, token string) (bool, error) {
	_, err := h.ParseVerifyToken(token)
	if nil != err {
		return false, err
	}

	return true, nil
}
