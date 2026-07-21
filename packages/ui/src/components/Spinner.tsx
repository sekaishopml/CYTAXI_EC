import { colors } from "@cytaxi/design-tokens";
import { cn } from "../utils/cn";

interface SpinnerProps {
  size?: number;
  color?: string;
  className?: string;
}

export function Spinner({ size = 20, color = colors.brand.green, className }: SpinnerProps) {
  return (
    <svg
      className={cn(className)}
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      style={{ animation: "spin 0.7s linear infinite" }}
    >
      <circle cx="12" cy="12" r="10" stroke={color} strokeWidth="3" opacity="0.15" />
      <path
        d="M12 2a10 10 0 019.95 9"
        stroke={color}
        strokeWidth="3"
        strokeLinecap="round"
      />
    </svg>
  );
}
