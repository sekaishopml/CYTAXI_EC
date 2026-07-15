import React from "react";

export function Footer() {
  return (
    <footer className="border-t border-border py-6 text-center text-sm text-muted">
      <div className="container mx-auto px-4 max-w-3xl">
        <p>CYTAXI &copy; {new Date().getFullYear()} &mdash; Conversational Mobility Platform</p>
      </div>
    </footer>
  );
}
