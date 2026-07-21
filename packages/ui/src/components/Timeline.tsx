import { colors, radius, spacing, typography, transition } from "@cytaxi/design-tokens";
import { cn } from "../utils/cn";

interface TimelineStep {
  id: string;
  label: string;
  description?: string;
  icon?: React.ReactNode;
  status: "pending" | "active" | "completed";
  eta?: string;
}

interface TimelineProps {
  steps: TimelineStep[];
  color?: string;
  className?: string;
}

export function Timeline({ steps, color = colors.brand.green, className }: TimelineProps) {
  return (
    <div className={cn(className)} style={{ display: "flex", flexDirection: "column", gap: 0 }}>
      {steps.map((step, i) => {
        const isLast = i === steps.length - 1;
        return (
          <div
            key={step.id}
            style={{
              display: "flex",
              gap: spacing[4],
              opacity: step.status === "pending" ? 0.4 : 1,
              transition: transition.fade,
            }}
          >
            {/* Timeline indicator */}
            <div style={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              width: 24,
              flexShrink: 0,
            }}>
              <div
                style={{
                  width: 12,
                  height: 12,
                  borderRadius: "50%",
                  background: step.status === "completed"
                    ? color
                    : step.status === "active"
                    ? color
                    : colors.border.medium,
                  boxShadow: step.status === "active"
                    ? `0 0 0 4px ${color}18`
                    : "none",
                  animation: step.status === "active"
                    ? "pulse 2s cubic-bezier(0.4,0,0.6,1) infinite"
                    : "none",
                  transition: transition.fade,
                  marginTop: 3,
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                }}
              >
                {step.status === "completed" && (
                  <svg width="8" height="8" viewBox="0 0 24 24" fill="white">
                    <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z" />
                  </svg>
                )}
              </div>
              {!isLast && (
                <div
                  style={{
                    width: 2,
                    flex: 1,
                    minHeight: 24,
                    background: step.status === "completed"
                      ? color
                      : colors.border.medium,
                    transition: transition.fade,
                    marginTop: 2,
                  }}
                />
              )}
            </div>

            {/* Content */}
            <div style={{
              flex: 1,
              paddingBottom: isLast ? 0 : spacing[6],
              display: "flex",
              justifyContent: "space-between",
              alignItems: "flex-start",
            }}>
              <div>
                <p style={{
                  margin: 0,
                  fontSize: typography.body.fontSize,
                  fontWeight: step.status === "active" ? 600 : 500,
                  color: step.status === "active" ? color : colors.text.primary,
                }}>
                  {step.label}
                </p>
                {step.description && (
                  <p style={{
                    margin: "2px 0 0",
                    fontSize: typography.caption.fontSize,
                    color: colors.text.muted,
                  }}>
                    {step.description}
                  </p>
                )}
              </div>
              {step.eta && (
                <span style={{
                  fontSize: typography.caption.fontSize,
                  fontWeight: 600,
                  color,
                  whiteSpace: "nowrap",
                  marginLeft: spacing[4],
                }}>
                  {step.eta}
                </span>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}
