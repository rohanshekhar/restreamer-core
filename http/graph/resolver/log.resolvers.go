package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"

	"github.com/datarhei/core/log"
)

func (r *queryResolver) Log(ctx context.Context) ([]string, error) {
	if r.LogBuffer == nil {
		r.LogBuffer = log.NewBufferWriter(log.Lsilent, 1)
	}

	events := r.LogBuffer.Events()

	formatter := log.NewConsoleFormatter(false)

	log := make([]string, len(events))

	for i, e := range events {
		log[i] = strings.TrimSpace(formatter.String(e))
	}

	return log, nil
}
