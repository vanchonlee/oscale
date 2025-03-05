// Package support contains the support calculation logic
package support

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/vanchonlee/oscale/internal/pkg/duration"
)

var (
	// ErrProvidersUtilisationSameLen is returned when the providers and max utilisation have different lengths
	ErrProvidersUtilisationSameLen = errors.New("providers and max utilisation should have the same length")
)

// GetUtilisation is a function that returns the utilisation for a specific offset and duration
type GetUtilisation func(ctx context.Context, offset, duration time.Duration) (float64, error)

// Entity defines a support entity
type Entity struct {
	Offset      duration.Duration `json:"offset"`
	Interval    duration.Duration `json:"interval"`
	Coefficient string            `json:"coefficient"`
}

// MustGetOffset returns the offset as a time.Duration
func (e *Entity) MustGetOffset() time.Duration {
	return e.Offset.MustDuration()
}

// MustGetInterval returns the interval as a time.Duration
func (e *Entity) MustGetInterval() time.Duration {
	return e.Interval.MustDuration()
}

// MustGetCoefficient returns the coefficient as a float64
func (e *Entity) MustGetCoefficient() float64 {
	coefficient, err := strconv.ParseFloat(e.Coefficient, 64)
	if err != nil {
		panic(err)
	}
	return coefficient
}

// Support defines a support
type Support struct {
	Entities []Entity `json:"entities,omitempty"`
}

// DeepCopyInto copies the Support object into the target Support object
func (s *Support) DeepCopyInto(out *Support) {
	out.Entities = make([]Entity, len(s.Entities))
	copy(out.Entities, s.Entities)
}

// Calculate calculates the support
func (s *Support) Calculate(ctx context.Context, providers []GetUtilisation, maxUtilisation []float64) (int, error) {
	l := log.FromContext(ctx)

	if len(providers) != len(maxUtilisation) {
		return 0, ErrProvidersUtilisationSameLen
	}

	support := 0
	for _, entity := range s.Entities {
		offset := entity.MustGetOffset()
		duration := entity.MustGetInterval()
		coefficient := entity.MustGetCoefficient()
		l.Info("Calculating support", "offset", offset, "duration", duration, "coefficient", coefficient)
		for i, provider := range providers {
			utilisation, err := provider(ctx, offset, duration)
			if err != nil {
				return 0, err
			}
			l.Info("Utilisation", "utilisation", utilisation)
			support = max(support, int(math.Ceil(utilisation*coefficient/maxUtilisation[i])))
		}

	}
	return support, nil
}
