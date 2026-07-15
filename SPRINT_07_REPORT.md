# Sprint 07 - Reporte Técnico

**Estado:** Listo para revisión

---

## Resumen

Sprint 07 (AI Orchestrator Foundation) completado. Se implementó el AI Orchestrator que decide cuándo utilizar un modelo LLM y cuándo resolver mediante reglas deterministas, con proveedores desacoplados y política de fallback.

---

## Archivos creados

| Archivo | Descripción |
|---------|-------------|
| `ai/orchestrator.go` | Orchestrator principal: decide LLM vs determinista, ejecuta providers |
| `ai/router.go` | Router: clasifica intención del usuario mediante classifiers |
| `ai/classifier.go` | RuleClassifier: clasificación por palabras clave |
| `ai/provider.go` | Provider interface + ProviderRegistry + mock provider |
| `ai/prompt.go` | PromptBuilder: construcción de system/user prompts |
| `ai/context.go` | ContextBuilder: contexto conversacional para el LLM |
| `ai/policy.go` | PolicyEvaluator: reglas por intent (UseLLM, min confidence) |
| `ai/fallback.go` | FallbackManager: retry, static response, fallback provider |
| `ai/observability.go` | MetricsCollector + AILogger para trazabilidad |
| `ai/README.md` | Documentación técnica del AI Orchestrator |

---

## Archivos modificados

Ninguno. Todos los archivos son nuevos dentro del Conversation Engine.

---

## Decisiones técnicas

1. **Router multiclassifier** — acepta múltiples classifiers en cadena. Primero en responder con confianza > 0.5 gana. Permite combinar reglas + LLM para clasificación.

2. **Provider desacoplado** — interfaz `Provider` con `Name()`, `Kind()`, `Complete()`. El `ProviderRegistry` selecciona el provider según capacidades requeridas. Cambiar de Qwen a DeepSeek requiere solo configurar otro provider.

3. **Política por intent** — cada intent tiene regla independiente: si usa LLM, confianza mínima y prioridad. `greeting` y `trip_status` son deterministas; `trip_request` y `support` usan LLM.

4. **Fallback progresivo** — 1) Reintento (2x), 2) Respuesta estática, 3) Provider alternativo. Cada nivel se activa solo si el anterior falla.

5. **Sin dependencias externas** — todos los providers son mock/interface. Las dependencias reales (openai, etc.) se agregarán en providers concretos.

6. **Métricas nativas** — `MetricsCollector` cuenta requests, decisiones, errores. Sin Prometheus dependency; exportable via interfaz.

---

## Riesgos

| Riesgo | Impacto | Mitigación |
|--------|---------|------------|
| Sin provider real implementado | Medio | Mock provider funcional para desarrollo |
| RuleClassifier limitado | Bajo | Diseñado para extenderse con classifiers LLM-based |
| Sin rate limiting | Bajo | Agregar en sprint futuro |
| Fallback no probado | Bajo | Estrategias definidas; probar con provider mock que falle |

---

## Mejoras futuras

- Implementar provider concreto para Qwen (Alibaba Cloud)
- Implementar provider concreto para DeepSeek
- Implementar provider concreto para OpenAI GPT
- Agregar classifier basado en LLM (usar LLM para clasificar intents)
- Agregar cache de respuestas del LLM
- Implementar rate limiting por sesión/proveedor
- Exportar métricas a Prometheus
- Agregar tracing distribuido (OpenTelemetry) para decisiones AI

---

## Siguiente Sprint recomendado

**Sprint 08 — AI Orchestrator: Provider Integration**

Integrar el AI Orchestrator con un proveedor LLM real:
- Implementar provider concreto (Qwen o DeepSeek)
- Conectar con MessagePipeline como MessageProcessor
- Probar flujo completo: mensaje → clasificación → LLM → respuesta
- Agregar tests de integración

---

## Definition of Done

- [x] AI Orchestrator creado
- [x] Proveedores desacoplados (Provider interface + Registry)
- [x] Sin lógica de negocio
- [x] Sin llamadas reales al LLM
- [x] Documentación incluida
- [x] Reporte entregado

---

*No se realizaron commits. No se realizó push. Esperando aprobación para continuar.*
