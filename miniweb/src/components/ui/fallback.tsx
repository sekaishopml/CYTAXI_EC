import React from "react";

export function Loading() {
  return (
    <div className="flex items-center justify-center py-12" role="status" aria-label="Loading">
      <div className="animate-spin rounded-full h-8 w-8 border-2 border-primary border-t-transparent" />
    </div>
  );
}

export function ErrorMessage({ message, onRetry }: { message: string; onRetry?: () => void }) {
  return (
    <div className="card text-center py-8" role="alert">
      <p className="text-danger font-medium mb-2">{message}</p>
      {onRetry && (
        <button onClick={onRetry} className="btn-secondary text-sm">Try again</button>
      )}
    </div>
  );
}
