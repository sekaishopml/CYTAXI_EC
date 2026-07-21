"use client";
import { Dispatch, SetStateAction } from "react";
import { Place } from "@/types";
import { colors, radius } from "@cytaxi/design-tokens";

interface FormStateProps {
  destQuery: string; setDestQuery: Dispatch<SetStateAction<string>>;
  destSuggestions: Place[];
  dest: Place | null;
  onSearch: (q: string) => void;
  onSelect: (place: Place) => void;
  onConfirm: () => void;
  onBack: () => void;
  onClearDest: () => void;
  loading: boolean;
  pickupAddress: string;
}

function getGreeting() {
  const h = new Date().getHours();
  if (h < 12) return "Buenos días";
  if (h < 18) return "Buenas tardes";
  return "Buenas noches";
}

export function FormState({ destQuery, setDestQuery, destSuggestions, dest, onSearch, onSelect, onConfirm, onBack, onClearDest, loading, pickupAddress }: FormStateProps) {
  const hasSelectedDest = dest !== null;

  return (
    <div style={{ display: "flex", flexDirection: "column", padding: "8px 18px 12px", gap: 10 }}>
      <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
        <button type="button" onClick={onBack} aria-label="Volver"
          style={{
            background: colors.surface.paperLight, border: "none", cursor: "pointer",
            width: 34, height: 34, borderRadius: radius.sm,
            display: "flex", alignItems: "center", justifyContent: "center",
            color: colors.textPrimary,
          }}>
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"><path d="M19 12H5m7-7-7 7 7 7"/></svg>
        </button>
        <p style={{ fontSize: 15, fontWeight: 600, color: colors.textPrimary, margin: 0, fontFamily: "'Space Grotesk', sans-serif", letterSpacing: "-0.01em" }}>
          {getGreeting()}
        </p>
      </div>

      <div style={{
        background: colors.surface.paperLight,
        borderRadius: radius.md, border: "1px solid rgba(0,0,0,0.06)",
      }}>
        {/* Origen — siempre tarjeta */}
        <div onClick={onBack} role="button" tabIndex={0} aria-label="Editar origen"
          onKeyDown={(e) => { if (e.key === "Enter" || e.key === " ") onBack(); }}
          style={{ display: "flex", alignItems: "center", gap: 10, padding: "12px 14px", cursor: "pointer" }}>
          <div style={{ width: 22, height: 22, borderRadius: "50%", background: colors.cobalt, flexShrink: 0, display: "flex", alignItems: "center", justifyContent: "center", color: "#fff", fontSize: 10, fontWeight: 700, fontFamily: "'Inter', sans-serif" }}>A</div>
          <div style={{ flex: 1, minWidth: 0 }}>
            <p style={{ fontSize: 9, fontWeight: 500, color: colors.textMuted, margin: 0, letterSpacing: "0.06em", textTransform: "uppercase", fontFamily: "'JetBrains Mono', monospace" }}>
              ORIGEN
            </p>
            <p style={{ fontSize: 13, fontWeight: 500, color: colors.textPrimary, margin: "2px 0 0", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" as const, fontFamily: "'Inter', sans-serif" }}>
              {pickupAddress}
            </p>
          </div>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke={colors.textMuted} strokeWidth="2" strokeLinecap="round"><path d="M9 18l6-6-6-6"/></svg>
        </div>

        <div style={{ height: 1, background: "rgba(0,0,0,0.06)", margin: "0 14px" }} />

        {hasSelectedDest ? (
          /* Destino — tarjeta seleccionada (no-editable) */
          <div onClick={onClearDest} role="button" tabIndex={0} aria-label="Cambiar destino"
            onKeyDown={(e) => { if (e.key === "Enter" || e.key === " ") onClearDest(); }}
            style={{
              display: "flex", alignItems: "center", gap: 10, padding: "12px 14px",
              cursor: "pointer", background: colors.cobaltBg,
              transition: "background 0.15s",
            }}>
            <div style={{
              width: 22, height: 22, borderRadius: "50%", background: colors.cobalt,
              flexShrink: 0, display: "flex", alignItems: "center", justifyContent: "center",
              color: "#fff", fontSize: 10, fontWeight: 700, fontFamily: "'Inter', sans-serif",
            }}>B</div>
            <div style={{ flex: 1, minWidth: 0 }}>
              <p style={{ fontSize: 9, fontWeight: 500, color: colors.cobalt, margin: 0, letterSpacing: "0.06em", textTransform: "uppercase", fontFamily: "'JetBrains Mono', monospace" }}>
                DESTINO
              </p>
              <p style={{ fontSize: 14, fontWeight: 600, color: colors.textPrimary, margin: "2px 0 0", fontFamily: "'Inter', sans-serif" }}>
                {dest?.name || destQuery}
              </p>
            </div>
            <span style={{ fontSize: 10, color: colors.textMuted, fontFamily: "'Inter', sans-serif", flexShrink: 0 }}>
              Cambiar
            </span>
          </div>
        ) : (
          /* Destino — input editable */
          <div style={{ display: "flex", alignItems: "center", gap: 10, padding: "12px 14px" }}>
            <div style={{ width: 22, height: 22, borderRadius: "50%", background: colors.cobaltLight, flexShrink: 0, display: "flex", alignItems: "center", justifyContent: "center", color: "#fff", fontSize: 10, fontWeight: 700, fontFamily: "'Inter', sans-serif" }}>B</div>
            <div style={{ flex: 1, minWidth: 0 }}>
              <p style={{ fontSize: 9, fontWeight: 500, color: colors.textMuted, margin: 0, letterSpacing: "0.06em", textTransform: "uppercase", fontFamily: "'JetBrains Mono', monospace" }}>
                DESTINO
              </p>
              <input
                aria-label="Buscar destino"
                placeholder="Buscar dirección o lugar"
                style={{
                  background: "transparent", border: "none", padding: 0, fontSize: 13,
                  fontFamily: "'Inter', sans-serif", color: colors.textPrimary, width: "100%",
                  outline: "none", fontWeight: 500, margin: "2px 0 0",
                }}
                value={destQuery}
                onChange={e => { setDestQuery(e.target.value); onSearch(e.target.value); }}
                autoFocus
              />
            </div>
            {destQuery && (
              <button type="button" onClick={onClearDest} aria-label="Limpiar destino"
                style={{
                  background: colors.surface.paper, border: "none", cursor: "pointer",
                  width: 26, height: 26, borderRadius: "50%",
                  display: "flex", alignItems: "center", justifyContent: "center",
                }}>
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke={colors.textMuted} strokeWidth="2.5" strokeLinecap="round"><path d="M18 6 6 18M6 6l12 12"/></svg>
              </button>
            )}
          </div>
        )}

        {destSuggestions.length > 0 && !hasSelectedDest && (
          <div style={{ borderTop: "1px solid rgba(0,0,0,0.05)" }}>
            {destSuggestions.slice(0, 4).map((p, i) => (
              <button type="button" key={i} onClick={() => onSelect(p)} aria-label={`Seleccionar ${p.name}`}
                style={{
                  padding: "10px 14px", display: "flex", alignItems: "center", gap: 10,
                  cursor: "pointer", border: "none", background: "transparent",
                  width: "100%", textAlign: "left", fontFamily: "'Inter', sans-serif",
                  borderBottom: i < Math.min(destSuggestions.length, 4) - 1 ? "1px solid rgba(0,0,0,0.04)" : "none",
                }}>
                <div style={{
                  width: 30, height: 30, borderRadius: radius.sm, background: colors.surface.paper,
                  display: "flex", alignItems: "center", justifyContent: "center", flexShrink: 0,
                }}>
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke={colors.textMuted} strokeWidth="2" strokeLinecap="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0118 0z"/><circle cx="12" cy="10" r="3"/></svg>
                </div>
                <div style={{ minWidth: 0, flex: 1 }}>
                  <p style={{ margin: 0, fontSize: 13, fontWeight: 500, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" as const, color: colors.textPrimary }}>{p.name}</p>
                  <p style={{ margin: "1px 0 0", fontSize: 11, color: colors.textMuted, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" as const }}>{p.address}</p>
                </div>
              </button>
            ))}
          </div>
        )}
      </div>

      <button type="button" onClick={onConfirm} disabled={loading}
        aria-label={hasSelectedDest ? "Buscar viaje" : "Selecciona tu destino en el mapa"}
        style={{
          width: "100%", height: 48, borderRadius: radius.xs, fontSize: 15, fontWeight: 600,
          fontFamily: "'Inter', sans-serif", border: hasSelectedDest ? "none" : "1px solid rgba(0,0,0,0.06)", cursor: "pointer",
          background: hasSelectedDest ? colors.cobalt : colors.surface.paperLight,
          color: hasSelectedDest ? "#fff" : colors.textPrimary,
          transition: "background 0.15s ease, opacity 0.15s ease",
          opacity: loading ? 0.6 : 1,
        }}>
        {loading ? (
          <span style={{ display: "inline-flex", alignItems: "center", gap: 8 }}>
            <span style={{ width: 16, height: 16, border: "2px solid rgba(255,255,255,0.3)", borderTopColor: "#fff", borderRadius: "50%", animation: "spin 0.7s linear infinite" }} />
            Calculando...
          </span>
        ) : hasSelectedDest ? "Buscar viaje" : "Selecciona tu destino"}
      </button>
    </div>
  );
}
