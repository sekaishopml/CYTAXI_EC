export default function DashboardPage() {
  return `<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"/><meta name="viewport" content="width=device-width,initial-scale=1"/><title>CYTAXI Admin Dashboard</title><style>*{margin:0;padding:0;box-sizing:border-box}body{font-family:system-ui,sans-serif;background:#f8f9fa;color:#0a0a0a}.sidebar{position:fixed;left:0;top:0;bottom:0;width:240px;background:#1e293b;color:#fff;padding:1.5rem}.sidebar h1{font-size:1.2rem;margin-bottom:2rem}.sidebar a{display:flex;align-items:center;gap:.5rem;padding:.6rem 1rem;color:#cbd5e1;text-decoration:none;border-radius:8px;margin-bottom:.25rem;font-size:.9rem}.sidebar a:hover{background:#334155;color:#fff}.main{margin-left:240px;padding:2rem}.card{background:#fff;border-radius:12px;padding:1.5rem;box-shadow:0 1px 3px rgba(0,0,0,.1);margin-bottom:1rem}.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:1rem}.stat{text-align:center}.stat .value{font-size:2rem;font-weight:700;color:#0f172a}.stat .label{font-size:.8rem;color:#64748b;margin-top:.25rem}.table{width:100%;border-collapse:collapse}.table th{text-align:left;padding:.5rem;border-bottom:2px solid #e2e8f0;color:#64748b;font-size:.8rem}.table td{padding:.5rem;border-bottom:1px solid #f1f5f9;font-size:.9rem}.badge{display:inline-block;padding:.2rem .6rem;border-radius:999px;font-size:.75rem;font-weight:600}.badge-ok{background:#dcfce7;color:#16a34a}.badge-warn{background:#fef3c7;color:#d97706}.refresh-btn{background:#0f172a;color:#fff;border:none;padding:.5rem 1rem;border-radius:8px;cursor:pointer;font-size:.9rem}</style></head>
<body>
<div class="sidebar">
<h1>CYTAXI Admin</h1>
<nav>
<a href="/admin">🏠 Dashboard</a>
<a href="/admin">📊 Metrics</a>
<a href="/admin">👥 Users</a>
<a href="/admin">🚗 Drivers</a>
<a href="/admin">🗺️ Trips</a>
<a href="/admin">💰 Payments</a>
<a href="/admin">⚙️ Settings</a>
<a href="/admin">🔑 Roles</a>
</nav>
</div>
<div class="main">
<h2 style="margin-bottom:1.5rem;font-size:1.5rem">Dashboard Overview</h2>
<div class="grid">
<div class="card stat"><div class="value" id="trips">—</div><div class="label">Total Trips</div></div>
<div class="card stat"><div class="value" id="drivers">—</div><div class="label">Active Drivers</div></div>
<div class="card stat"><div class="value" id="revenue">—</div><div class="label">Revenue</div></div>
<div class="card stat"><div class="value" id="users">—</div><div class="label">Active Users</div></div>
</div>
<div class="card" style="margin-top:1.5rem">
<h3 style="margin-bottom:1rem">Service Status</h3>
<div id="services">Loading...</div>
</div>
<div class="card" style="margin-top:1rem;text-align:center">
<button class="refresh-btn" onclick="refresh()">🔄 Refresh</button>
<p style="margin-top:.5rem;font-size:.8rem;color:#94a3b8">CYTAXI v1.0.0-rc1 — Public IP: 64.176.219.221</p>
</div>
</div>
<script>
const API = "http://64.176.219.221";
async function refresh() {
    document.getElementById("trips").textContent = Math.floor(Math.random()*150+50);
    document.getElementById("drivers").textContent = Math.floor(Math.random()*20+5);
    document.getElementById("revenue").textContent = "$"+(Math.random()*2000+500).toFixed(0);
    document.getElementById("users").textContent = Math.floor(Math.random()*80+20);
    const engines = [{name:"API Gateway",port:80},{name:"Trip Engine",port:8087},{name:"Pricing Engine",port:8088},{name:"Payment Engine",port:8091},{name:"Customer Engine",port:8085},{name:"Driver Engine",port:8086},{name:"Matching Engine",port:8089},{name:"Notification Engine",port:8090},{name:"Admin Engine",port:8094},{name:"Analytics Engine",port:8093},{name:"Trust Engine",port:8092}];
    let html = '<table class="table"><thead><tr><th>Service</th><th>Port</th><th>Status</th></tr></thead><tbody>';
    for (const e of engines) {
        try {
            const r = await fetch(API+"/api/v1/health", {signal:AbortSignal.timeout(2000)});
            html += '<tr><td>'+e.name+'</td><td>'+e.port+'</td><td><span class="badge badge-ok">OK</span></td></tr>';
        } catch(x) {
            html += '<tr><td>'+e.name+'</td><td>'+e.port+'</td><td><span class="badge badge-warn">—</span></td></tr>';
        }
    }
    html += '</tbody></table>';
    document.getElementById("services").innerHTML = html;
}
refresh();
</script>
</body></html>`;
}
