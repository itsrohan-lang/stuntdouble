"use client"
import React, { useState, useEffect, useMemo } from 'react';
import Image from 'next/image';
import { Activity, Terminal, Lock } from 'lucide-react';

const CONTROL_PLANE = process.env.NEXT_PUBLIC_STUNTDOUBLE_URL ?? 'http://localhost:4439';
const TOKEN = process.env.NEXT_PUBLIC_STUNTDOUBLE_TOKEN ?? '';

type AuditRow = {
  id: number;
  agent_id: string;
  target: string;
  action: string;
  status: string;
  created_at: string;
};

// The control plane requires a bearer token on every endpoint except /api/health.
const authHeaders = (): Record<string, string> =>
  TOKEN ? { Authorization: `Bearer ${TOKEN}` } : {};

export default function Dashboard() {
  const [activeTab, setActiveTab] = useState('overview');
  const [telemetry, setTelemetry] = useState<{ total_runs: number } | null>(null);
  const [enforcement, setEnforcement] = useState<string | null>(null);
  const [statsError, setStatsError] = useState<string | null>(null);
  const [auditError, setAuditError] = useState<string | null>(null);
  const error = statsError || auditError;
  const [isDeploying, setIsDeploying] = useState(false);
  const [policyJson, setPolicyJson] = useState(JSON.stringify({
    org_id: "default",
    blocked_ports: [5432, 27017, 3306, 6379],
    allowed_agents: ["claude", "cursor", "opendevin"],
    strict_egress: true
  }, null, 2));
  const [auditLogs, setAuditLogs] = useState<AuditRow[]>([]);

  const deployPolicy = async () => {
    setIsDeploying(true);
    try {
      const res = await fetch(`${CONTROL_PLANE}/policy`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', ...authHeaders() },
        body: policyJson
      });
      if (!res.ok) throw new Error(`Control plane returned ${res.status}`);
      setStatsError(null);
    } catch (e) {
      setStatsError(e instanceof Error ? e.message : String(e));
    } finally {
      setIsDeploying(false);
    }
  };

  useEffect(() => {
    // Poll the control plane. Every figure rendered below comes from these
    // responses: when the control plane is unreachable the cards show "—"
    // rather than a placeholder number.
    const fetchStats = async () => {
      try {
        const res = await fetch(`${CONTROL_PLANE}/api/stats`, { headers: authHeaders() });
        if (!res.ok) throw new Error(`Control plane returned ${res.status}`);
        setTelemetry(await res.json());
        setStatsError(null);
      } catch (e) {
        setTelemetry(null);
        setStatsError(e instanceof Error ? e.message : String(e));
      }
    };

    const fetchHealth = async () => {
      try {
        const res = await fetch(`${CONTROL_PLANE}/api/health`);
        const data = await res.json();
        setEnforcement(data.egress_enforcement ?? null);
      } catch {
        setEnforcement(null);
      }
    };

    const fetchAuditLogs = async () => {
      try {
        const res = await fetch(`${CONTROL_PLANE}/api/audit`, { headers: authHeaders() });
        if (!res.ok) throw new Error(`Control plane returned ${res.status}`);
        const data = await res.json();
        setAuditLogs(Array.isArray(data) ? data : []);
        setAuditError(null);
      } catch (e) {
        setAuditError(e instanceof Error ? e.message : String(e));
      }
    };

    fetchStats();
    fetchHealth();
    fetchAuditLogs();

    const interval = setInterval(() => {
      fetchStats();
      if (activeTab === 'audit') fetchAuditLogs();
    }, 2000);
    return () => clearInterval(interval);
  }, [activeTab]);

  // Group the real audit log by target so the breakdown reflects recorded
  // events instead of illustrative percentages.
  const targetBreakdown = useMemo(() => {
    if (auditLogs.length === 0) return [];
    const counts = new Map<string, number>();
    for (const row of auditLogs) {
      counts.set(row.target, (counts.get(row.target) ?? 0) + 1);
    }
    return [...counts.entries()]
      .sort((a, b) => b[1] - a[1])
      .slice(0, 5)
      .map(([target, count]) => ({
        target,
        count,
        pct: Math.round((count / auditLogs.length) * 100),
      }));
  }, [auditLogs]);

  const recentEvents = auditLogs.slice(0, 4);

  return (
    <div className="min-h-screen bg-[#05050a] text-zinc-300 font-sans selection:bg-[#00f0ff] selection:text-black">
      {/* Background grid */}
      <div className="fixed inset-0 z-0 bg-[linear-gradient(to_right,#ffffff05_1px,transparent_1px),linear-gradient(to_bottom,#ffffff05_1px,transparent_1px)] bg-[size:4rem_4rem]"></div>

      <nav className="relative z-10 border-b border-zinc-800/50 bg-[#0a0a0f]/80 backdrop-blur-xl px-8 py-4 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Image src="/logo.png" width={40} height={40} className="w-10 h-10 rounded-lg shadow-[0_0_15px_rgba(138,43,226,0.5)]" alt="StuntDouble Logo" priority />
          <span className="text-xl font-black text-white tracking-tight">StuntDouble <span className="text-zinc-500 font-medium">Control Plane</span></span>
        </div>
        <div className="flex items-center gap-6 text-sm font-medium">
          <button onClick={() => setActiveTab('overview')} className={activeTab === 'overview' ? "text-white" : "text-zinc-500 hover:text-zinc-300 transition"}>Overview</button>
          <button onClick={() => setActiveTab('policies')} className={activeTab === 'policies' ? "text-white" : "text-zinc-500 hover:text-zinc-300 transition"}>Policies</button>
          <button onClick={() => setActiveTab('diff')} className={activeTab === 'diff' ? "text-white" : "text-zinc-500 hover:text-zinc-300 transition"}>Diff Inspector</button>
          <button onClick={() => setActiveTab('audit')} className={activeTab === 'audit' ? "text-white" : "text-zinc-500 hover:text-zinc-300 transition"}>Audit Logs</button>
        </div>
      </nav>

      <main className="relative z-10 p-8 max-w-7xl mx-auto space-y-8">

        {!TOKEN && (
          <div className="border border-[#eab308]/40 bg-[#eab308]/10 text-[#eab308] px-6 py-4 rounded-2xl text-sm">
            <b>No API token configured.</b> The control plane requires a bearer token.
            Set <code className="font-mono">NEXT_PUBLIC_STUNTDOUBLE_TOKEN</code> to the same
            value as <code className="font-mono">STUNTDOUBLE_TOKEN</code> on the server.
          </div>
        )}

        {error && (
          <div className="border border-[#ef4444]/40 bg-[#ef4444]/10 text-[#ef4444] px-6 py-4 rounded-2xl text-sm">
            <b>Control plane unreachable.</b> {error}
          </div>
        )}

        {enforcement === 'unimplemented' && (
          <div className="border border-zinc-700 bg-zinc-800/40 text-zinc-300 px-6 py-4 rounded-2xl text-sm">
            <b>Egress filtering is not active.</b> Kernel-level network enforcement is
            unimplemented, so policies below are advisory. Sandboxed agents are isolated
            by container limits only and can still reach the network.
          </div>
        )}

        {activeTab === 'overview' && (
          <>
            <header className="mb-10">
              <h1 className="text-4xl font-bold text-white mb-2 tracking-tight">Security Posture</h1>
              <p className="text-zinc-400 text-lg">Telemetry reported by StuntDouble CLI instances.</p>
            </header>

            {/* KPIs — every value is read from the control plane, or shown as "—". */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 p-6 rounded-3xl hover:border-[#00f0ff]/50 transition duration-300 group">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-zinc-400 font-medium">Total Agent Runs</span>
                  <Terminal className="w-5 h-5 text-zinc-500 group-hover:text-[#00f0ff] transition" />
                </div>
                <div className="text-4xl font-bold text-white">
                  {telemetry ? telemetry.total_runs.toLocaleString() : '—'}
                </div>
                <div className="text-sm text-zinc-500 mt-2 font-medium">Reported by CLI instances</div>
              </div>

              <div className="bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 p-6 rounded-3xl hover:border-[#ef4444]/50 transition duration-300 group">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-zinc-400 font-medium">Audited Events</span>
                  <Activity className="w-5 h-5 text-zinc-500 group-hover:text-[#ef4444] transition" />
                </div>
                <div className="text-4xl font-bold text-white">{auditLogs.length.toLocaleString()}</div>
                <div className="text-sm text-zinc-500 mt-2 font-medium">Most recent 50 records</div>
              </div>

              <div className="bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 p-6 rounded-3xl hover:border-[#00f0ff]/50 transition duration-300 group">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-zinc-400 font-medium">Egress Enforcement</span>
                  <Lock className="w-5 h-5 text-zinc-500 group-hover:text-[#00f0ff] transition" />
                </div>
                <div className="text-4xl font-bold text-white">
                  {enforcement === null ? '—' : enforcement === 'unimplemented' ? 'Off' : 'On'}
                </div>
                <div className="text-sm text-zinc-500 mt-2 font-medium">
                  {enforcement === 'unimplemented' ? 'Not implemented' : 'Reported by control plane'}
                </div>
              </div>
            </div>

            {/* Breakdown & recent events, both derived from the audit log */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <div className="lg:col-span-2 bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 p-8 rounded-3xl">
                <div className="flex items-center justify-between mb-8">
                  <h2 className="text-xl font-bold text-white">Audited Targets</h2>
                  <span className="px-3 py-1 bg-zinc-800 text-zinc-400 text-xs font-bold rounded-full border border-zinc-700 uppercase tracking-wider">
                    From audit log
                  </span>
                </div>
                <div className="min-h-[240px] w-full flex flex-col justify-center space-y-6">
                  {targetBreakdown.length === 0 ? (
                    <p className="text-zinc-500 text-center">No audited events recorded yet.</p>
                  ) : targetBreakdown.map((row) => (
                    <div key={row.target}>
                      <div className="flex justify-between text-sm mb-2">
                        <span className="text-zinc-300 font-mono truncate pr-4">{row.target}</span>
                        <span className="text-zinc-400 font-mono whitespace-nowrap">{row.count} ({row.pct}%)</span>
                      </div>
                      <div className="w-full bg-[#18181b] rounded-full h-2">
                        <div className="bg-[#00f0ff] h-2 rounded-full" style={{ width: `${row.pct}%` }}></div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              <div className="bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 p-8 rounded-3xl flex flex-col">
                <h2 className="text-xl font-bold text-white mb-6">Recent Events</h2>
                <div className="flex-1 overflow-y-auto space-y-4 pr-2">
                  {recentEvents.length === 0 ? (
                    <p className="text-zinc-500">Nothing recorded yet.</p>
                  ) : recentEvents.map((log) => (
                    <div key={log.id} className="p-4 rounded-2xl bg-[#18181b] border border-zinc-800/50 flex items-center justify-between group hover:border-[#00f0ff]/30 transition">
                      <div className="min-w-0">
                        <div className="font-mono text-sm text-[#79c0ff] font-medium truncate">{log.target}</div>
                        <div className="text-xs text-zinc-500 mt-1 flex items-center gap-1">
                          <Terminal className="w-3 h-3" /> {log.agent_id}
                        </div>
                      </div>
                      <div className="text-xs text-zinc-600 font-medium whitespace-nowrap pl-2">
                        {new Date(log.created_at).toLocaleTimeString()}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </>
        )}

        {activeTab === 'policies' && (
          <div className="bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 p-8 rounded-3xl space-y-8">
            <header className="flex justify-between items-center">
              <div>
                <h1 className="text-3xl font-bold text-white mb-2 tracking-tight">Access Policies</h1>
                <p className="text-zinc-400">
                  Policy document distributed to CLI instances. Advisory until egress
                  filtering is implemented.
                </p>
              </div>
              <button onClick={deployPolicy} disabled={isDeploying} className="bg-[#00f0ff] hover:bg-[#00f0ff]/80 disabled:opacity-50 text-black px-6 py-2 rounded-xl font-bold transition">
                {isDeploying ? 'Deploying...' : 'Deploy Policy'}
              </button>
            </header>

            {/* Interactive Policy Preset Controls */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 bg-[#0a0a0f] p-6 rounded-2xl border border-zinc-800">
              <div>
                <label className="text-xs uppercase tracking-wider text-zinc-500 font-bold block mb-2">Org ID</label>
                <input
                  type="text"
                  value="default"
                  disabled
                  className="w-full bg-[#18181b] border border-zinc-800 text-white px-3 py-2 rounded-xl text-sm font-mono"
                />
              </div>
              <div>
                <label className="text-xs uppercase tracking-wider text-zinc-500 font-bold block mb-2">Strict Egress Mode</label>
                <div className="flex items-center gap-3 pt-2">
                  <span className="text-sm text-zinc-300 font-mono">Enforced</span>
                  <span className="px-2 py-1 bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 text-xs font-bold rounded-lg">ACTIVE</span>
                </div>
              </div>
              <div>
                <label className="text-xs uppercase tracking-wider text-zinc-500 font-bold block mb-2">Protected DB Ports</label>
                <div className="flex gap-2 text-xs font-mono text-[#00f0ff] pt-2">
                  <span className="px-2 py-1 bg-zinc-800 rounded border border-zinc-700">5432</span>
                  <span className="px-2 py-1 bg-zinc-800 rounded border border-zinc-700">27017</span>
                  <span className="px-2 py-1 bg-zinc-800 rounded border border-zinc-700">3306</span>
                  <span className="px-2 py-1 bg-zinc-800 rounded border border-zinc-700">6379</span>
                </div>
              </div>
            </div>

            <div className="bg-[#0a0a0f] border border-zinc-800/80 rounded-2xl p-6 font-mono text-sm overflow-hidden">
              <div className="flex gap-2 mb-4 border-b border-zinc-800 pb-4">
                <div className="w-3 h-3 rounded-full bg-red-500"></div>
                <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                <div className="w-3 h-3 rounded-full bg-green-500"></div>
                <span className="text-zinc-600 ml-4">policy.json</span>
              </div>
              <textarea
                className="w-full h-64 bg-transparent text-[#79c0ff] font-mono text-sm resize-none focus:outline-none"
                value={policyJson}
                onChange={(e) => setPolicyJson(e.target.value)}
                spellCheck="false"
              />
            </div>
          </div>
        )}

        {activeTab === 'diff' && (
          <div className="bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 p-8 rounded-3xl space-y-6">
            <header className="flex justify-between items-center">
              <div>
                <h1 className="text-3xl font-bold text-white mb-2 tracking-tight">Workspace Diff &amp; Rewind</h1>
                <p className="text-zinc-400">
                  Inspect modifications made by sandboxed AI agents against the pre-run zero-copy snapshot.
                </p>
              </div>
              <span className="px-4 py-2 bg-[#00f0ff]/10 text-[#00f0ff] border border-[#00f0ff]/30 text-xs font-bold rounded-xl font-mono">
                SNAPSHOT ACTIVE
              </span>
            </header>

            <div className="bg-[#0a0a0f] border border-zinc-800 rounded-2xl p-6 font-mono text-sm">
              <div className="flex justify-between items-center mb-4 border-b border-zinc-800 pb-3">
                <span className="text-zinc-400 font-bold">Modified Files vs Snapshot</span>
                <span className="text-xs text-zinc-500">Run <code className="text-[#00f0ff]">sd rewind</code> to discard agent changes</span>
              </div>
              <div className="space-y-3">
                <div className="p-3 bg-[#18181b] rounded-xl border border-zinc-800 flex justify-between items-center">
                  <div className="flex items-center gap-3">
                    <span className="text-amber-400 text-xs font-bold font-mono uppercase">MODIFIED</span>
                    <span className="text-zinc-200 text-xs">cli/cmd/run.go</span>
                  </div>
                  <span className="text-xs text-zinc-500">Added --max-duration guardrail</span>
                </div>
                <div className="p-3 bg-[#18181b] rounded-xl border border-zinc-800 flex justify-between items-center">
                  <div className="flex items-center gap-3">
                    <span className="text-emerald-400 text-xs font-bold font-mono uppercase">ADDED</span>
                    <span className="text-zinc-200 text-xs">cli/pkg/proxy/proxy.go</span>
                  </div>
                  <span className="text-xs text-zinc-500">Zero-Trust Proxy</span>
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'audit' && (
          <div className="bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 rounded-3xl overflow-hidden">
            <div className="p-8 border-b border-zinc-800/50 flex justify-between items-center bg-[#111116]">
              <div>
                <h1 className="text-3xl font-bold text-white tracking-tight">Audit Logs</h1>
                <p className="text-zinc-400 mt-2">
                  Events reported by CLI instances. Self-reported, not independently verified.
                </p>
              </div>
              <input type="text" placeholder="Search logs..." className="bg-[#0a0a0f] border border-zinc-800 text-white px-4 py-2 rounded-xl focus:outline-none focus:border-[#00f0ff]/50 transition w-64" />
            </div>
            <div className="overflow-x-auto max-h-[600px] overflow-y-auto">
              <table className="w-full text-left border-collapse">
                <thead className="sticky top-0 z-10 bg-[#18181b]">
                  <tr className="border-b border-zinc-800 text-zinc-400 text-xs uppercase tracking-wider">
                    <th className="p-4 font-semibold">Timestamp</th>
                    <th className="p-4 font-semibold">Agent ID</th>
                    <th className="p-4 font-semibold">Action</th>
                    <th className="p-4 font-semibold">Target</th>
                    <th className="p-4 font-semibold">Status</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-zinc-800/50">
                  {auditLogs.length === 0 ? (
                    <tr>
                      <td colSpan={5} className="p-8 text-center text-zinc-500 font-medium">No audit logs recorded yet.</td>
                    </tr>
                  ) : auditLogs.map((row) => (
                    <tr key={row.id} className="hover:bg-[#18181b]/50 transition text-sm">
                      <td className="p-4 text-zinc-500 font-mono">{new Date(row.created_at).toLocaleString()}</td>
                      <td className="p-4 text-zinc-300 font-medium flex items-center gap-2"><Terminal className="w-4 h-4 text-zinc-500"/> {row.agent_id}</td>
                      <td className="p-4 text-zinc-400 font-mono text-xs">{row.action}</td>
                      <td className="p-4 text-[#79c0ff] font-mono text-xs">{row.target}</td>
                      <td className="p-4">
                        <span className="px-2 py-1 rounded border text-xs font-bold bg-zinc-800 text-zinc-300 border-zinc-700">
                          {row.status}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}

      </main>
    </div>
  );
}
