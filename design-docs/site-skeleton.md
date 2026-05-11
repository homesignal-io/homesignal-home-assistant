import React, { useState } from "react";

const pages = [
  "Login",
  "Dashboard",
  "Sites",
  "Add Site",
  "Site Detail",
  "Provision Device",
  "Device Detail",
  "Alerts",
  "Company Settings",
  "Team",
];

export default function HomeSignalPortalPreview() {
  const [page, setPage] = useState("Dashboard");

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900 flex">
      {page !== "Login" && (
        <aside className="w-64 bg-white border-r border-slate-200 p-5">
          <div className="text-xl font-semibold mb-8">HomeSignal</div>
          <nav className="space-y-1">
            {pages.filter(p => p !== "Login").map(p => (
              <button
                key={p}
                onClick={() => setPage(p)}
                className={`w-full text-left px-3 py-2 rounded-lg text-sm ${
                  page === p
                    ? "bg-slate-900 text-white"
                    : "hover:bg-slate-100"
                }`}
              >
                {p}
              </button>
            ))}
          </nav>
        </aside>
      )}

      <main className="flex-1 p-8">
        {page === "Login" && <Login setPage={setPage} />}
        {page === "Dashboard" && <Dashboard setPage={setPage} />}
        {page === "Sites" && <Sites setPage={setPage} />}
        {page === "Add Site" && <AddSite />}
        {page === "Site Detail" && <SiteDetail setPage={setPage} />}
        {page === "Provision Device" && <ProvisionDevice />}
        {page === "Device Detail" && <DeviceDetail />}
        {page === "Alerts" && <Alerts />}
        {page === "Company Settings" && <CompanySettings />}
        {page === "Team" && <Team />}
      </main>
    </div>
  );
}

function Card({ title, children, action }) {
  return (
    <section className="bg-white rounded-2xl border border-slate-200 p-6 shadow-sm">
      <div className="flex justify-between items-start mb-4">
        <h2 className="text-lg font-semibold">{title}</h2>
        {action}
      </div>
      {children}
    </section>
  );
}

function Button({ children, onClick, variant = "primary" }) {
  const classes =
    variant === "secondary"
      ? "border border-slate-300 bg-white hover:bg-slate-50"
      : "bg-slate-900 text-white hover:bg-slate-700";

  return (
    <button onClick={onClick} className={`px-4 py-2 rounded-lg text-sm ${classes}`}>
      {children}
    </button>
  );
}

function Login({ setPage }) {
  return (
    <div className="max-w-md mx-auto mt-24 bg-white border rounded-2xl p-8 shadow-sm">
      <h1 className="text-2xl font-semibold mb-2">HomeSignal</h1>
      <p className="text-slate-500 mb-6">Sign in to your integrator account.</p>
      <Label>Email</Label>
      <Input />
      <Label>Password</Label>
      <Input type="password" />
      <Button onClick={() => setPage("Dashboard")}>Sign in</Button>
      <div className="mt-4 text-sm text-slate-500 flex gap-4">
        <span>Forgot password?</span>
        <span>Request access</span>
      </div>
    </div>
  );
}

function Dashboard({ setPage }) {
  return (
    <Page title="Dashboard">
      <div className="grid grid-cols-4 gap-4 mb-6">
        <Stat label="Online sites" value="18" />
        <Stat label="Offline sites" value="1" />
        <Stat label="Needs attention" value="3" />
        <Stat label="Pending updates" value="7" />
      </div>

      <Card title="Attention needed" action={<Button onClick={() => setPage("Add Site")}>Add site</Button>}>
        <List items={[
          "Smith Residence — Device offline — Last seen 42 minutes ago",
          "Lee Residence — Backup failed — Last attempt yesterday",
          "Patel Residence — Home Assistant update available",
        ]} />
      </Card>
    </Page>
  );
}

function Sites({ setPage }) {
  return (
    <Page title="Sites">
      <div className="flex gap-3 mb-6">
        <Input placeholder="Search sites" />
        <Button onClick={() => setPage("Add Site")}>Add site</Button>
      </div>

      <div className="grid gap-4">
        {["Smith Residence", "Lee Residence", "Patel Residence"].map((site, i) => (
          <button
            key={site}
            onClick={() => setPage("Site Detail")}
            className="bg-white border rounded-2xl p-5 text-left hover:border-slate-400"
          >
            <div className="font-semibold">{site}</div>
            <div className="text-sm text-slate-500">
              Status: {i === 1 ? "Needs attention" : "Online"} · Devices: 1 · Last seen: Now
            </div>
          </button>
        ))}
      </div>
    </Page>
  );
}

function AddSite() {
  return (
    <Page title="Add site">
      <Card title="Site details">
        <Label>Customer / site name</Label>
        <Input />
        <Label>Address</Label>
        <Input />
        <Label>Customer contact</Label>
        <Input />
        <Label>Notes</Label>
        <textarea className="w-full border rounded-lg p-3 mb-4" rows="4" />
        <Label>Deployment template</Label>
        <select className="w-full border rounded-lg p-3 mb-6">
          <option>Standard Home Assistant Deployment v1</option>
        </select>
        <Button>Create site</Button>
      </Card>
    </Page>
  );
}

function SiteDetail({ setPage }) {
  return (
    <Page title="Smith Residence">
      <div className="grid grid-cols-2 gap-6">
        <Card title="Health">
          <List items={[
            "Home Assistant: Healthy",
            "HomeSignal add-on: Connected",
            "Backups: Last successful backup yesterday",
            "Updates: 2 available",
            "Alerts: None",
          ]} />
        </Card>

        <Card title="Actions">
          <div className="flex flex-wrap gap-3">
            <Button onClick={() => setPage("Provision Device")}>Add Home Assistant device</Button>
            <Button variant="secondary">Trigger backup</Button>
            <Button variant="secondary" onClick={() => setPage("Device Detail")}>View device</Button>
            <Button variant="secondary">Release device</Button>
          </div>
        </Card>
      </div>
    </Page>
  );
}

function ProvisionDevice() {
  return (
    <Page title="Add Home Assistant device">
      <Card title="Provisioning session">
        <p className="text-sm text-slate-500 mb-6">Site: Smith Residence</p>
        <Steps items={[
          "Install the HomeSignal Manager add-on in Home Assistant.",
          "Open the HomeSignal Manager Web UI inside Home Assistant.",
          "Enter the pairing code shown by the add-on.",
          "Confirm bind.",
        ]} />
        <div className="flex gap-3 mt-6">
          <Button>Install on Home Assistant</Button>
          <Button variant="secondary">Regenerate session</Button>
        </div>
        <div className="mt-6">
          <Label>Pairing code</Label>
          <Input placeholder="___ ___" />
          <Button>Pair device</Button>
        </div>
        <p className="text-sm text-slate-500 mt-4">Status: Waiting for pairing · Expires in: 14:32</p>
      </Card>
    </Page>
  );
}

function DeviceDetail() {
  return (
    <Page title="Home Assistant Green">
      <div className="grid grid-cols-2 gap-6">
        <Card title="Device details">
          <List items={[
            "Status: Online",
            "Site: Smith Residence",
            "Hardware: Home Assistant Green",
            "Home Assistant version: 2026.5.x",
            "Supervisor: Healthy",
            "HomeSignal add-on: 1.0.0",
            "Connection: Connected",
            "Storage: Healthy",
          ]} />
        </Card>
        <Card title="Actions">
          <div className="flex flex-wrap gap-3">
            <Button>Request diagnostics</Button>
            <Button variant="secondary">Restart add-on</Button>
            <Button variant="secondary">Trigger backup</Button>
            <Button variant="secondary">Release device</Button>
          </div>
        </Card>
      </div>
    </Page>
  );
}

function Alerts() {
  return (
    <Page title="Alerts">
      <Card title="Active alerts">
        <List items={[
          "Smith Residence — High — Device offline — First seen today 3:14 PM",
          "Lee Residence — Medium — Backup failed — First seen yesterday",
        ]} />
        <div className="flex gap-3 mt-6">
          <Button>Acknowledge</Button>
          <Button variant="secondary">Resolve</Button>
        </div>
      </Card>
    </Page>
  );
}

function CompanySettings() {
  return (
    <Page title="Company settings">
      <Card title="Company profile">
        <Label>Company name</Label>
        <Input defaultValue="ABC Smart Homes" />
        <Label>Support email</Label>
        <Input defaultValue="support@example.com" />
        <Label>Support phone</Label>
        <Input />
        <Label>Default timezone</Label>
        <Input defaultValue="America/New_York" />
        <Button>Save changes</Button>
      </Card>
    </Page>
  );
}

function Team() {
  return (
    <Page title="Team">
      <Card title="Users">
        <List items={[
          "Jamie Smith — Owner — jamie@example.com",
          "Alex Lee — Technician — alex@example.com",
        ]} />
      </Card>

      <div className="mt-6">
        <Card title="Invite user">
          <Label>Email</Label>
          <Input />
          <Label>Role</Label>
          <select className="w-full border rounded-lg p-3 mb-6">
            <option>Technician</option>
            <option>Admin</option>
            <option>Read-only</option>
          </select>
          <Button>Send invite</Button>
        </Card>
      </div>
    </Page>
  );
}

function Page({ title, children }) {
  return (
    <>
      <h1 className="text-3xl font-semibold mb-8">{title}</h1>
      {children}
    </>
  );
}

function Stat({ label, value }) {
  return (
    <div className="bg-white border rounded-2xl p-5 shadow-sm">
      <div className="text-3xl font-semibold">{value}</div>
      <div className="text-sm text-slate-500 mt-1">{label}</div>
    </div>
  );
}

function List({ items }) {
  return (
    <ul className="space-y-3 text-sm">
      {items.map(item => (
        <li key={item} className="border-b last:border-b-0 pb-3 last:pb-0">
          {item}
        </li>
      ))}
    </ul>
  );
}

function Steps({ items }) {
  return (
    <ol className="space-y-3">
      {items.map((item, i) => (
        <li key={item} className="flex gap-3">
          <span className="w-7 h-7 rounded-full bg-slate-900 text-white flex items-center justify-center text-sm">
            {i + 1}
          </span>
          <span>{item}</span>
        </li>
      ))}
    </ol>
  );
}

function Label({ children }) {
  return <label className="block text-sm font-medium mb-1 mt-4">{children}</label>;
}

function Input(props) {
  return <input {...props} className="w-full border rounded-lg p-3 mb-4" />;
}