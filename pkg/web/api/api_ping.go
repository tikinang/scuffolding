package api

import "context"

type PingIn struct {
	Ping bool `json:"ping" binding:"required"`
}

type PingOut struct {
	Pong bool `json:"pong" binding:"required"`
}

func (r *Handler) Ping(ctx context.Context, in PingIn) (PingOut, error) {
	return PingOut{Pong: true}, nil
}
