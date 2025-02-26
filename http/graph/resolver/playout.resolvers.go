package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/datarhei/core/http/graph/models"
	"github.com/datarhei/core/playout"
)

func (r *queryResolver) PlayoutStatus(ctx context.Context, id string, input string) (*models.RawAVstream, error) {
	addr, err := r.Restream.GetPlayout(id, input)
	if err != nil {
		return nil, fmt.Errorf("unknown process or input: %w", err)
	}

	path := "/v1/status"

	data, err := r.playoutRequest(http.MethodGet, addr, path, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	status := playout.Status{}

	if err := json.Unmarshal(data, &status); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s := &models.RawAVstream{}
	s.UnmarshalPlayout(status)

	return s, nil
}
