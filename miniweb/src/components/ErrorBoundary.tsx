"use client";
import { Component, ErrorInfo, ReactNode } from "react";
import { colors, radius } from "@cytaxi/design-tokens";

interface ErrorBoundaryProps {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error, info: ErrorInfo) => void;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    this.props.onError?.(error, info);
    try {
      fetch("/api/v1/telemetry", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          events: [{ event: "REACT_ERROR", error: error.message, stack: error.stack, timestamp: Date.now() }],
          source: "miniweb",
        }),
        keepalive: true,
      });
    } catch {}
  }

  render() {
    if (this.state.hasError) {
      if (this.props.fallback) return this.props.fallback;
      return (
        <div style={{
          padding: 40, textAlign: "center", display: "flex",
          flexDirection: "column", alignItems: "center", gap: 16,
        }}>
          <div style={{ fontSize: 48, marginBottom: 8 }}>⚠️</div>
          <p style={{ fontSize: 18, fontWeight: 600, color: colors.textPrimary, margin: 0 }}>
            Algo salió mal
          </p>
          <p style={{ fontSize: 13, color: colors.textMuted, margin: 0, maxWidth: 260 }}>
            Hubo un error inesperado. Puedes intentar recargar la página.
          </p>
          <button type="button"
            onClick={() => window.location.reload()}
            style={{
              padding: "12px 28px", borderRadius: radius.md, fontSize: 14, fontWeight: 600,
              border: "none", cursor: "pointer", background: colors.cobalt,
              color: "#fff", fontFamily: "Inter",
            }}
          >
            Recargar página
          </button>
          {this.state.error && (
            <details style={{ fontSize: 11, color: colors.textMuted, maxWidth: 300, textAlign: "left" }}>
              <summary style={{ cursor: "pointer" }}>Detalles técnicos</summary>
              <pre style={{ marginTop: 8, fontSize: 10, whiteSpace: "pre-wrap", wordBreak: "break-all" }}>
                {this.state.error.message}
              </pre>
            </details>
          )}
        </div>
      );
    }
    return this.props.children;
  }
}
