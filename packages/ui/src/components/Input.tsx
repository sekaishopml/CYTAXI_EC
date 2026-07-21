import { useRef, useState } from "react";
import { colors, radius, spacing, typography, fontFamily, transition, duration } from "@cytaxi/design-tokens";
import { cn } from "../utils/cn";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  icon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  containerStyle?: React.CSSProperties;
}

export function Input({
  label,
  error,
  icon,
  rightIcon,
  containerStyle,
  className,
  style,
  ...rest
}: InputProps) {
  const [focused, setFocused] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: spacing[1], ...containerStyle }}>
      {label && (
        <label
          onClick={() => inputRef.current?.focus()}
          style={{
            fontSize: typography.caption.fontSize,
            fontWeight: typography.label.fontWeight,
            color: colors.text.muted,
            letterSpacing: "0.03em",
          }}
        >
          {label}
        </label>
      )}
      <div
        onClick={() => inputRef.current?.focus()}
        style={{
          display: "flex",
          alignItems: "center",
          gap: spacing[3],
          height: 48,
          padding: `0 ${spacing[4]}px`,
          borderRadius: radius.md,
          border: `1.5px solid ${error ? colors.status.error : focused ? colors.brand.green : colors.border.medium}`,
          background: colors.surface.white,
          transition: `border-color ${duration.fade}ms ease, box-shadow ${duration.fade}ms ease`,
          boxShadow: focused ? `0 0 0 3px ${colors.brand.greenBorder}` : "none",
          cursor: "text",
          ...style,
        }}
      >
        {icon && <span style={{ color: colors.text.muted, display: "flex" }}>{icon}</span>}
        <input
          ref={inputRef}
          className={cn(className)}
          onFocus={() => setFocused(true)}
          onBlur={() => setFocused(false)}
          style={{
            flex: 1,
            border: "none",
            outline: "none",
            background: "transparent",
            fontSize: typography.body.fontSize,
            fontWeight: typography.body.fontWeight,
            fontFamily: fontFamily.sans,
            color: colors.text.primary,
            width: "100%",
          }}
          {...rest}
        />
        {rightIcon && <span style={{ color: colors.text.muted, display: "flex" }}>{rightIcon}</span>}
      </div>
      {error && (
        <span style={{ fontSize: typography.caption.fontSize, color: colors.status.error, marginTop: 2 }}>
          {error}
        </span>
      )}
    </div>
  );
}
