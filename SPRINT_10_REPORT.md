# Sprint 10 - Reporte Técnico

**Estado:** Listo para revisión

---

## Resumen

Sprint 10 (Mobility Decision Engine) completado. Se creó la infraestructura del Engine responsable de coordinar las decisiones de despacho, con pipeline modular, estrategias de asignación intercambiables y contratos desacoplados de Conversation, Geospatial y Policy Engines.

---

## Archivos creados

| Archivo | Descripción |
|---------|-------------|
| `go.mod` | Módulo Go del Engine |
| `cmd/mobility/main.go` | Bootstrap + health + 3 estrategias registradas |
| `domain/assignment/assignment.go` | Assignment entity con estados (pending→proposed→accepted/rejected/expired/completed/cancelled) |
| `domain/candidate/candidate.go` | Candidate, CandidateSet, Location, Vehicle |
| `domain/decision/decision.go` | Decision, DecisionContext, StepSummary |
| `domain/decision/errors.go` | Errores tipados del dominio |
| `application/port/port.go` | DispatcherInputPort, CandidateInputPort, AssignmentInputPort |
| `application/dispatcher/dispatcher.go` | DispatcherCoordinator: orquesta el flujo completo |
| `application/pipeline/pipeline.go` | DecisionPipeline: steps chain + strategy selection |
| `application/pipeline/builder.go` | CandidateBuilder: filters + scorers, ProximityFilter, AvailabilityFilter |
| `infrastructure/strategy/strategy.go` | StrategyRegistry + NearestDriver + HighestRated + BalancedScore |
| `events/definition.go` | 13 eventos del mobility engine |
| `config/config.go` | Config (port, default strategy, max candidates, timeout) |
| `README.md` | Documentación completa |
| `Dockerfile` | Dockerfile multi-stage |

---

## Archivos modificados

| Archivo | Cambio |
|---------|--------|
| `go.work` | Se agregó `./backend/engines/mobility` |

---

## Arquitectura aplicada

```
DispatcherCoordinator.Dispatch(ctx, DecisionContext)
       ↓
CandidateFinder.FindCandidates()
       ↓ (CandidateSet)
DecisionPipeline.Execute()
  ├── PipelineStep 1: ProximityFilter
  ├── PipelineStep 2: AvailabilityFilter
  ├── PipelineStep 3: VehicleTypeFilter
  └── (extensible)
       ↓ (filtered CandidateSet)
Strategy.Select() ← StrategyRegistry
  ├── NearestDriver     (menor distancia)
  ├── HighestRated      (mayor score)
  └── BalancedScore     (pesos configurables)
       ↓
Decision (status, selected_driver, strategy, score, pipeline_summary)
```

**Pipeline extensible:** cualquier `PipelineStep` puede agregarse. El pipeline ejecuta en orden y cada paso recibe el set de candidatos filtrado por el paso anterior.

**Estrategias intercambiables:** `StrategyRegistry` permite registrar y seleccionar estrategias por nombre en tiempo de ejecución.

---

## Decisiones técnicas

1. **Pipeline + Strategy separation** — los pasos del pipeline filtran/modifican candidatos; la estrategia selecciona el ganador. Separación clara de responsabilidades.
2. **StrategyRegistry** — permite cambiar algoritmo de asignación sin modificar el núcleo del Engine. Las estrategias se registran en bootstrap y se seleccionan por nombre en configuración.
3. **DecisionContext como transporte** — contiene toda la información del viaje, usuario, origen, destino, requisitos y datos de pricing/policy. Se pasa a través de todo el pipeline.
4. **CandidateBuilder** — provee filters y scorers reutilizables para construir el CandidateSet desde drivers crudos.
5. **BalancedScore strategy** — combinación ponderada de distancia normalizada y score normalizado, con pesos configurables.
6. **Eventos granulares** — 13 eventos que cubren cada etapa del dispatch (start, candidates, pipeline steps, strategy, assignment).

---

## Riesgos

| Riesgo | Impacto | Mitigación |
|--------|---------|------------|
| CandidateFinder no implementado | Alto | Interfaz definida; depende de Driver Engine |
| Sin algoritmo de asignación real | Medio | 3 estrategias implementadas (proximidad, rating, balance) |
| Sin persistencia | Bajo | Repositorios se agregarán en sprint futuro |

---

## Deuda técnica

- CandidateFinder es nil en bootstrap — requiere Driver Engine para implementar
- PipelineSteps actuales son no-operativos (devuelven mismos candidatos)
- Sin tests de integración del pipeline completo

---

## Mejoras futuras

- Implementar CandidateFinder real conectado a Driver Engine + Geospatial Engine
- Agregar pipeline step de validación de políticas (Policy Engine)
- Agregar pipeline step de pricing (Revenue Engine)
- Agregar timeout real en dispatcher
- Implementar cola de despacho para alta concurrencia
- Agregar endpoints REST para consultar estado de dispatch

---

## Commit sugerido

```
feat(mobility): create Mobility Decision Engine infrastructure
```

---

## Definition of Done

- [x] Engine creado
- [x] Pipeline definido (DecisionPipeline con steps + strategy)
- [x] Interfaces completas (DispatcherInputPort, CandidateInputPort, AssignmentInputPort)
- [x] Sin algoritmo de asignación complejo
- [x] Sin IA
- [x] Sin integraciones externas
- [x] Documentación completa

---

*No se realizaron commits. No se realizó push. Esperando aprobación para continuar.*
