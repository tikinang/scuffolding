package site

import "context"

func (r *Handler) Index(ctx context.Context) (map[string]any, error) {
	return map[string]any{"title": "super title"}, nil
}
