package saga

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sekaishopml/cytaxi/backend/integration/contracts"
	"github.com/sekaishopml/cytaxi/backend/integration/eventbus"
)

type Step struct {
	Name       string
	Execute    func(ctx context.Context, data any) (any, error)
	Compensate func(ctx context.Context, data any) error
}

type SagaDefinition struct {
	ID        string
	Name      string
	Steps     []Step
	Timeout   time.Duration
	Retries   int
}

type SagaCoordinator struct {
	definitions map[string]SagaDefinition
	bus         eventbus.Bus
	logger      *slog.Logger
}

func NewSagaCoordinator(bus eventbus.Bus, logger *slog.Logger) *SagaCoordinator {
	return &SagaCoordinator{
		definitions: make(map[string]SagaDefinition),
		bus:         bus,
		logger:      logger,
	}
}

func (sc *SagaCoordinator) Register(def SagaDefinition) {
	sc.definitions[def.ID] = def
}

func (sc *SagaCoordinator) Execute(ctx context.Context, sagaID string, data any) error {
	def, ok := sc.definitions[sagaID]
	if !ok {
		return fmt.Errorf("saga %s not found", sagaID)
	}

	sc.logger.Info("saga started", "saga", def.Name, "steps", len(def.Steps))

	var result any = data
	for i, step := range def.Steps {
		sc.logger.Info("saga step executing", "saga", def.Name, "step", step.Name, "index", i+1)
		out, err := step.Execute(ctx, result)
		if err != nil {
			sc.logger.Error("saga step failed", "saga", def.Name, "step", step.Name, "error", err)
			go sc.compensate(ctx, def, i)
			return fmt.Errorf("saga %s step %s: %w", def.Name, step.Name, err)
		}
		result = out
	}

	sc.logger.Info("saga completed", "saga", def.Name)
	return nil
}

func (sc *SagaCoordinator) compensate(ctx context.Context, def SagaDefinition, failedIndex int) {
	sc.logger.Info("saga compensating", "saga", def.Name, "from_step", failedIndex)
	for i := failedIndex; i >= 0; i-- {
		step := def.Steps[i]
		if step.Compensate != nil {
			sc.logger.Info("saga compensating step", "step", step.Name)
			step.Compensate(ctx, nil)
		}
	}
}
