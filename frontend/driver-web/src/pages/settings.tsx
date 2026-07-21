import React from "react";

export default function SettingsPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Settings</h1>
      <div className="card space-y-4">
        <div className="flex justify-between items-center"><span>Notifications</span><input type="checkbox" defaultChecked /></div>
        <div className="flex justify-between items-center"><span>Auto-accept trips</span><input type="checkbox" /></div>
        <div className="flex justify-between items-center"><span>Language</span><select className="input w-auto"><option>English</option><option>Español</option></select></div>
      </div>
    </div>
  );
}
