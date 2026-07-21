"use client";

interface BottomSheetHeaderProps {
  title?: string;
  subtitle?: string;
}

export function BottomSheetHeader({ title, subtitle }: BottomSheetHeaderProps) {
  if (!title && !subtitle) return null;
  return (
    <div style={{ padding: "0 4px 12px" }}>
      {title && (
        <p style={{
          margin: 0, fontSize: 17, fontWeight: 600,
          fontFamily: "'Space Grotesk', sans-serif",
          color: "#121212",
        }}>
          {title}
        </p>
      )}
      {subtitle && (
        <p style={{
          margin: "2px 0 0", fontSize: 13,
          fontFamily: "'Inter', sans-serif",
          color: "rgba(0,0,0,0.45)",
        }}>
          {subtitle}
        </p>
      )}
    </div>
  );
}
