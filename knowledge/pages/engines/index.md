# Engines (Motores de Negocio)

Los 10 motores definidos en CYDIGITAL-BLUEPRINT representan capacidades de negocio, no unidades de despliegue. Un motor puede implementarse como uno o más microservicios.

## 1. Conversation Engine
- **Capítulo Blueprint**: 22
- **Propósito**: Procesamiento de lenguaje natural, reconocimiento de intención
- **Canal primario**: WhatsApp
- **Responsabilidades**: Entender solicitudes en lenguaje natural, extraer ubicaciones, resolver direcciones ambiguas
- **Implementación actual**: ❌ No implementado

## 2. Geospatial Intelligence Engine
- **Capítulo Blueprint**: 23
- **Propósito**: Geocoding, rutas, tráfico, ETA
- **Responsabilidades**: Geocoding directo e inverso, cálculo de rutas, optimización
- **Implementación actual**: ✅ `backend/engines/geospatial/` — Go, OSRM + Nominatim
- **API**: `/v1/geocode`, `/v1/reverse`, `/v1/routes`, `/v1/eta`

## 3. Mobility Decision Engine
- **Capítulo Blueprint**: 24
- **Propósito**: Dispatch inteligente, asignación de conductores
- **Responsabilidades**: Matching conductor-viaje, cola de prioridad, ETA ranking
- **Implementación actual**: ❌ No implementado

## 4. Revenue Intelligence Engine
- **Capítulo Blueprint**: 25
- **Propósito**: Pricing determinístico, tarifas, promociones
- **Responsabilidades**: Cálculo de tarifas, surge pricing, descuentos, redondeo
- **Implementación actual**: ❌ No implementado
- **Regla**: Siempre determinístico — la IA nunca decide precios

## 5. Customer Intelligence Engine
- **Capítulo Blueprint**: 26
- **Propósito**: Memoria del cliente, recomendaciones
- **Responsabilidades**: Perfil, favoritos, historial, recomendaciones personalizadas
- **Implementación actual**: ❌ No implementado

## 6. Driver Intelligence Engine
- **Capítulo Blueprint**: 27
- **Propósito**: Reputación, excelencia operacional
- **Responsabilidades**: Rating de conductores, historial, wallet, documentos
- **Implementación actual**: ❌ No implementado

## 7. Trust Intelligence Platform
- **Capítulo Blueprint**: 28
- **Propósito**: Riesgo, fraude, seguridad
- **Responsabilidades**: Detección de fraude, verificación, auditoría
- **Implementación actual**: ❌ No implementado

## 8. Notification Engine
- **Capítulo Blueprint**: 9.14
- **Propósito**: Notificaciones multicanal
- **Responsabilidades**: WhatsApp, push, SMS, email
- **Implementación actual**: ❌ No implementado

## 9. Analytics Engine
- **Capítulo Blueprint**: 9.17
- **Propósito**: Métricas, forecasting
- **Responsabilidades**: Dashboards, KPIs, reportes históricos
- **Implementación actual**: ❌ No implementado

## 10. Policy Engine
- **Capítulo Blueprint**: 16
- **Propósito**: Reglas de negocio, autorización
- **Responsabilidades**: Evaluación de políticas, permisos, Zero Trust
- **Implementación actual**: ❌ No implementado
