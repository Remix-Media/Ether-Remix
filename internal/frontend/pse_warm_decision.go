package frontend

import (
	"context"
	"time"

	"github.com/braedonsaunders/ether/internal/core"
)

// WarmDecision is what the warm-up of a stream says about its sources.
type WarmDecision = core.WarmDecision

const (
	WarmUnknown   = core.WarmUnknown
	WarmHasSource = core.WarmHasSource
	WarmNoSource  = core.WarmNoSource
)

// LazyStreamDecider is optionally implemented by LazyStreams. After Warm, it
// waits at most budget for the verdict of the resolver for the same request.
// Anything but WarmNoSource means "keep offering the stream".
type LazyStreamDecider interface {
	WarmDecision(ctx context.Context, credential string, req PlanRequest, budget time.Duration) WarmDecision
}

// NoReleaseMemo is optionally implemented by the playback manager: it tells
// whether the warm-up of a title for this viewer and device recently found
// that no release exists (internal/playback/pse_warm_decision.go).
type NoReleaseMemo interface {
	PSENoReleaseRecent(user *User, key RatingKey, caps ClientCaps) bool
}
