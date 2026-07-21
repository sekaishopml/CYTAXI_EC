import { colors, fontWeight } from "@cytaxi/design-tokens";
import { cn } from "../utils/cn";

interface AvatarProps {
  src?: string;
  name?: string;
  size?: number;
  rating?: number;
  className?: string;
  style?: React.CSSProperties;
}

export function Avatar({ src, name, size = 48, rating, className, style }: AvatarProps) {
  const initials = name
    ? name.split(" ").map(n => n[0]).join("").slice(0, 2).toUpperCase()
    : "?";

  return (
    <div
      className={cn(className)}
      style={{
        position: "relative",
        width: size,
        height: size,
        borderRadius: "50%",
        flexShrink: 0,
        overflow: "hidden",
        background: src ? "none" : colors.brand.greenBg,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        border: `2px solid ${colors.surface.white}`,
        boxShadow: `0 0 0 1px ${colors.border.medium}`,
        ...style,
      }}
    >
      {src ? (
        <img
          src={src}
          alt={name || "Avatar"}
          style={{ width: "100%", height: "100%", objectFit: "cover" }}
        />
      ) : (
        <span style={{
          fontSize: size * 0.4,
          fontWeight: fontWeight.semibold,
          color: colors.brand.green,
        }}>
          {initials}
        </span>
      )}
      {rating !== undefined && (
        <span
          style={{
            position: "absolute",
            bottom: -2,
            right: -2,
            background: colors.surface.white,
            borderRadius: "50%",
            padding: 2,
            boxShadow: "0 1px 4px rgba(0,0,0,0.1)",
            fontSize: size * 0.2,
            fontWeight: 700,
            lineHeight: 1,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          ⭐{rating.toFixed(1)}
        </span>
      )}
    </div>
  );
}
