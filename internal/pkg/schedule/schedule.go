// Package schedule provides the logic to parse and evaluate schedules
package schedule

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// EntityError wraps errors that occur during parsing of a schedule
type EntityError struct {
	Err error
}

// Error returns the error message
func (e EntityError) Error() string {
	if e.Err == nil {
		return "failed to parse schedule"
	}
	return errors.Wrap(e.Err, "failed to parse schedule").Error()
}

// Is checks if the target error is an EntityError
func (e EntityError) Is(target error) bool {
	var scheduleError EntityError
	ok := errors.As(target, &scheduleError)
	return ok
}

// Entity defines a schedule entity
type Entity struct {
	CronStart string `json:"cronStart"`
	CronEnd   string `json:"cronEnd"`
	Count     int32  `json:"count"`
}

// Schedule defines a schedule
type Schedule struct {
	Entities []Entity `json:"entities,omitempty"`
}

// DeepCopyInto copies the Schedule object into the target Schedule object
func (s *Schedule) DeepCopyInto(out *Schedule) {
	out.Entities = make([]Entity, len(s.Entities))
	copy(out.Entities, s.Entities)
}

// GetCountAt returns the minimum count at a given time
func (s *Schedule) GetCountAt(ctx context.Context, at time.Time) (int, error) {
	minimum := 0
	logger := zap.L()

	for _, schedule := range s.Entities {
		shouldScale, err := schedule.Check(at)
		if err != nil {
			logger.Error("Failed to Check schedule",
				zap.Error(err))
			return -1, err
		}
		if shouldScale {
			logger.Info("Schedule triggered",
				zap.Int("amount", int(schedule.Count)),
				zap.String("cronUp", schedule.CronStart),
				zap.String("cronDown", schedule.CronEnd),
			)
			minimum = int(schedule.Count)
		} else {
			logger.Info("Schedule not triggered",
				zap.Int("amount", int(schedule.Count)),
				zap.String("cronUp", schedule.CronStart),
				zap.String("cronDown", schedule.CronEnd),
			)
		}
	}
	return minimum, nil
}

// Check checks if the schedule is active at a given time
func (e *Entity) Check(at time.Time) (bool, error) {
	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	// Parse the cronUp expression
	cronUpSchedule, err := cronParser.Parse(e.CronStart)
	if err != nil {
		return false, EntityError{Err: err}
	}

	// Parse the cronDown expression
	cronDownSchedule, err := cronParser.Parse(e.CronEnd)
	if err != nil {
		return false, EntityError{Err: err}
	}

	// Get the next scheduled times for cronUp and cronDown
	nextUp := cronUpSchedule.Next(at)
	nextDown := cronDownSchedule.Next(at)

	return nextDown.Before(nextUp), nil
}
