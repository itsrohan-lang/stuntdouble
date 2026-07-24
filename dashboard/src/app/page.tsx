"use client"
import React, { useState, useEffect } from 'react';
import { Activity, ShieldAlert, Users, ServerCrash, Terminal, Lock, Globe } from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, AreaChart, Area } from 'recharts';

export default function Dashboard() {
  const [activeTab, setActiveTab] = useState('overview');
  const [telemetry, setTelemetry] = useState({ total_runs: 0, blocked_commands: 0 });
  const [isDeploying, setIsDeploying] = useState(false);
  const [policyJson, setPolicyJson] = useState(JSON.stringify({
    org_id: "ent_global_updated",
    blocked_ports: [5432, 27017, 3306, 6379, 8080],
    allowed_agents: ["claude", "cursor", "opendevin"],
    strict_egress: true
  }, null, 2));
  const [auditLogs, setAuditLogs] = useState([]);
  const [isSimulating, setIsSimulating] = useState(false);
  
  const deployPolicy = async () => {
    setIsDeploying(true);
    try {
      await fetch('http://localhost:4439/policy', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: policyJson
      });
      setTimeout(() => setIsDeploying(false), 1000);
    } catch (e) {
      console.error(e);
      setIsDeploying(false);
    }
  };

  const fetchAuditLogs = async () => {
    try {
      const res = await fetch('http://localhost:4439/api/audit');
      const data = await res.json();
      if (data) setAuditLogs(data);
    } catch (e) {
      console.error("Failed to fetch audit logs", e);
    }
  };

  const simulateRogueAttack = async () => {
    setIsSimulating(true);
    try {
      await fetch('http://localhost:4439/api/audit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          agent_id: "rogue-agent-007",
          target: "api.stripe.com/v1/payouts",
          action: "NETWORK_EXFILTRATION",
          status: "Blocked (eBPF Kernel)",
          created_at: new Date().toISOString()
        })
      });
      await fetchAuditLogs();
      setTimeout(() => setIsSimulating(false), 500);
    } catch (e) {
      console.error(e);
      setIsSimulating(false);
    }
  };

  useEffect(() => {
    // Poll the Control Plane API for live stats every 2 seconds
    const fetchStats = async () => {
      try {
        const res = await fetch('http://localhost:4439/api/stats');
        const data = await res.json();
        // If data isn't initialized on backend, keep existing state
        if (data.total_runs > 0 || data.blocked_commands > 0) {
          setTelemetry({ total_runs: data.total_runs, blocked_commands: data.blocked_commands });
        }
      } catch (e) {
        console.error("Failed to fetch live stats", e);
      }
    };
    
    // Initial fetch
    fetchStats();
    fetchAuditLogs();
    
    const interval = setInterval(() => {
      fetchStats();
      if (activeTab === 'audit') fetchAuditLogs();
    }, 2000);
    return () => clearInterval(interval);
  }, [activeTab]);

  // Simulated live data feed
  const [data, setData] = useState([
    { time: '10:00', blocked: 12 },
    { time: '10:05', blocked: 19 },
    { time: '10:10', blocked: 15 },
    { time: '10:15', blocked: 25 },
    { time: '10:20', blocked: 22 },
    { time: '10:25', blocked: 30 },
  ]);

  useEffect(() => {
    // In a real app, this would fetch from http://localhost:8080/metrics
    setTelemetry({ total_runs: 1402, blocked_commands: 843 });
  }, []);

  return (
    <div className="min-h-screen bg-[#05050a] text-zinc-300 font-sans selection:bg-[#00f0ff] selection:text-black">
      {/* Background grid */}
      <div className="fixed inset-0 z-0 bg-[linear-gradient(to_right,#ffffff05_1px,transparent_1px),linear-gradient(to_bottom,#ffffff05_1px,transparent_1px)] bg-[size:4rem_4rem]"></div>

      <nav className="relative z-10 border-b border-zinc-800/50 bg-[#0a0a0f]/80 backdrop-blur-xl px-8 py-4 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <img src="/logo.png" className="w-10 h-10 rounded-lg shadow-[0_0_15px_rgba(138,43,226,0.5)]" alt="StuntDouble Logo" />
          <span className="text-xl font-black text-white tracking-tight">StuntDouble <span className="text-zinc-500 font-medium">Control Plane</span></span>
        </div>
        <div className="flex items-center gap-6 text-sm font-medium">
          <button onClick={() => setActiveTab('overview')} className={activeTab === 'overview' ? "text-white" : "text-zinc-500 hover:text-zinc-300 transition"}>Overview</button>
          <button onClick={() => setActiveTab('policies')} className={activeTab === 'policies' ? "text-white" : "text-zinc-500 hover:text-zinc-300 transition"}>Policies</button>
          <button onClick={() => setActiveTab('audit')} className={activeTab === 'audit' ? "text-white" : "text-zinc-500 hover:text-zinc-300 transition"}>Audit Logs</button>
          <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-[#8a2be2] to-[#00f0ff] flex items-center justify-center text-white font-bold ml-4 shadow-[0_0_15px_rgba(0,240,255,0.3)]">
            CTO
          </div>
        </div>
      </nav>

      <main className="relative z-10 p-8 max-w-7xl mx-auto space-y-8">
        
        {activeTab === 'overview' && (
          <>
            <header className="mb-10">
              <h1 className="text-4xl font-bold text-white mb-2 tracking-tight">Global Security Posture</h1>
              <p className="text-zinc-400 text-lg">Real-time telemetry across all organizational AI sandboxes.</p>
            </header>

            {/* KPIs */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
              <div className="bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 p-6 rounded-3xl hover:border-[#00f0ff]/50 transition duration-300 group">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-zinc-400 font-medium">Total Agent Runs</span>
                  <Terminal className="w-5 h-5 text-zinc-500 group-hover:text-[#00f0ff] transition" />
                </div>
                <div className="text-4xl font-bold text-white">{telemetry.total_runs.toLocaleString()}</div>
                <div className="text-sm text-[#00f0ff] mt-2 font-medium">+14% from last week</div>
              </div>
              
              <div className="bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 p-6 rounded-3xl hover:border-[#ef4444]/50 transition duration-300 group">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-zinc-400 font-medium">Blocked Requests</span>
                  <ShieldAlert className="w-5 h-5 text-zinc-500 group-hover:text-[#ef4444] transition" />
                </div>
                <div className="text-4xl font-bold text-white">{telemetry.blocked_commands.toLocaleString()}</div>
                <div className="text-sm text-[#ef4444] mt-2 font-medium">12 critical severity</div>
              </div>

              <div className="bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 p-6 rounded-3xl hover:border-[#8a2be2]/50 transition duration-300 group">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-zinc-400 font-medium">Active Sandboxes</span>
                  <Activity className="w-5 h-5 text-zinc-500 group-hover:text-[#8a2be2] transition" />
                </div>
                <div className="text-4xl font-bold text-white">42</div>
                <div className="text-sm text-zinc-500 mt-2 font-medium">Across 3 clusters</div>
              </div>

              <div className="bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 p-6 rounded-3xl hover:border-[#00f0ff]/50 transition duration-300 group">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-zinc-400 font-medium">Enforcement Mode</span>
                  <Lock className="w-5 h-5 text-zinc-500 group-hover:text-[#00f0ff] transition" />
                </div>
                <div className="text-4xl font-bold text-white">Strict</div>
                <div className="text-sm text-[#00f0ff] mt-2 font-medium">eBPF Blocking Enabled</div>
              </div>
            </div>

            {/* Charts & Logs */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <div className="lg:col-span-2 bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 p-8 rounded-3xl">
                <div className="flex items-center justify-between mb-8">
                  <h2 className="text-xl font-bold text-white">Threat Vector Analysis</h2>
                  <span className="px-3 py-1 bg-[#ef4444]/10 text-[#ef4444] text-xs font-bold rounded-full border border-[#ef4444]/20 uppercase tracking-wider">Live Blocks</span>
                </div>
                <div className="h-[300px] w-full flex flex-col justify-center space-y-6">
                  <div>
                    <div className="flex justify-between text-sm mb-2"><span className="text-zinc-300">Database Ports (5432, 27017)</span><span className="text-[#ef4444] font-mono">42%</span></div>
                    <div className="w-full bg-[#18181b] rounded-full h-2"><div className="bg-[#ef4444] h-2 rounded-full" style={{width: '42%'}}></div></div>
                  </div>
                  <div>
                    <div className="flex justify-between text-sm mb-2"><span className="text-zinc-300">Cloud Metadata APIs (169.254.x.x)</span><span className="text-[#f97316] font-mono">28%</span></div>
                    <div className="w-full bg-[#18181b] rounded-full h-2"><div className="bg-[#f97316] h-2 rounded-full" style={{width: '28%'}}></div></div>
                  </div>
                  <div>
                    <div className="flex justify-between text-sm mb-2"><span className="text-zinc-300">Stripe/Payment APIs</span><span className="text-[#eab308] font-mono">18%</span></div>
                    <div className="w-full bg-[#18181b] rounded-full h-2"><div className="bg-[#eab308] h-2 rounded-full" style={{width: '18%'}}></div></div>
                  </div>
                  <div>
                    <div className="flex justify-between text-sm mb-2"><span className="text-zinc-300">Internal K8s DNS</span><span className="text-[#3b82f6] font-mono">12%</span></div>
                    <div className="w-full bg-[#18181b] rounded-full h-2"><div className="bg-[#3b82f6] h-2 rounded-full" style={{width: '12%'}}></div></div>
                  </div>
                </div>
              </div>

              <div className="bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 p-8 rounded-3xl flex flex-col">
                <h2 className="text-xl font-bold text-white mb-6">Recent Violations</h2>
                <div className="flex-1 overflow-y-auto space-y-4 pr-2">
                  {[
                    { target: 'api.stripe.com', agent: 'claude-code', time: '2m ago' },
                    { target: 'postgres:5432', agent: 'opendevin', time: '14m ago' },
                    { target: 'github.com/private', agent: 'cursor', time: '1h ago' },
                    { target: 's3.amazonaws.com', agent: 'aider', time: '3h ago' },
                  ].map((log, i) => (
                    <div key={i} className="p-4 rounded-2xl bg-[#18181b] border border-zinc-800/50 flex items-center justify-between group hover:border-[#ef4444]/30 transition">
                      <div>
                        <div className="font-mono text-sm text-[#ef4444] font-medium">{log.target}</div>
                        <div className="text-xs text-zinc-500 mt-1 flex items-center gap-1">
                          <Terminal className="w-3 h-3" /> {log.agent}
                        </div>
                      </div>
                      <div className="text-xs text-zinc-600 font-medium">{log.time}</div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </>
        )}

        {activeTab === 'policies' && (
          <div className="bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 p-8 rounded-3xl">
            <header className="mb-8 flex justify-between items-center">
              <div>
                <h1 className="text-3xl font-bold text-white mb-2 tracking-tight">Access Policies</h1>
                <p className="text-zinc-400">Manage Zero-Trust rules distributed to all StuntDouble eBPF nodes.</p>
              </div>
              <button onClick={deployPolicy} disabled={isDeploying} className="bg-[#00f0ff] hover:bg-[#00f0ff]/80 disabled:opacity-50 text-black px-6 py-2 rounded-xl font-bold transition">
                {isDeploying ? 'Deploying...' : 'Deploy Policy'}
              </button>
            </header>
            
            <div className="bg-[#0a0a0f] border border-zinc-800/80 rounded-2xl p-6 font-mono text-sm overflow-hidden">
              <div className="flex gap-2 mb-4 border-b border-zinc-800 pb-4">
                <div className="w-3 h-3 rounded-full bg-red-500"></div>
                <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                <div className="w-3 h-3 rounded-full bg-green-500"></div>
                <span className="text-zinc-600 ml-4">.stuntdouble.yaml</span>
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

        {activeTab === 'audit' && (
          <div className="bg-[#111116]/80 backdrop-blur-md border border-zinc-800/50 rounded-3xl overflow-hidden shadow-[0_0_50px_rgba(239,68,68,0.05)]">
            <div className="p-8 border-b border-zinc-800/50 flex justify-between items-center bg-[#111116]">
              <div>
                <h1 className="text-3xl font-bold text-white tracking-tight">Audit Logs</h1>
                <p className="text-zinc-400 mt-2">Immutable enterprise ledger of all agent actions.</p>
              </div>
              <div className="flex gap-4">
                <button 
                  onClick={simulateRogueAttack}
                  disabled={isSimulating}
                  className="bg-[#ef4444]/10 hover:bg-[#ef4444]/20 border border-[#ef4444]/50 text-[#ef4444] px-4 py-2 rounded-xl font-bold transition flex items-center gap-2"
                >
                  <ShieldAlert className="w-4 h-4" />
                  {isSimulating ? 'Simulating...' : 'Simulate Rogue Attack'}
                </button>
                <input type="text" placeholder="Search logs..." className="bg-[#0a0a0f] border border-zinc-800 text-white px-4 py-2 rounded-xl focus:outline-none focus:border-[#00f0ff]/50 transition w-64" />
              </div>
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
                      <td colSpan={5} className="p-8 text-center text-zinc-500 font-medium">No audit logs recorded yet. Try simulating an attack!</td>
                    </tr>
                  ) : auditLogs.map((row: any) => (
                    <tr key={row.id} className={`hover:bg-[#18181b]/50 transition text-sm ${row.action === 'NETWORK_EXFILTRATION' ? 'bg-[#ef4444]/5' : ''}`}>
                      <td className="p-4 text-zinc-500 font-mono">{new Date(row.created_at).toLocaleString()}</td>
                      <td className="p-4 text-zinc-300 font-medium flex items-center gap-2"><Terminal className="w-4 h-4 text-zinc-500"/> {row.agent_id}</td>
                      <td className="p-4 text-zinc-400 font-mono text-xs">{row.action}</td>
                      <td className="p-4 text-[#79c0ff] font-mono text-xs">{row.target}</td>
                      <td className="p-4">
                        <span className={`px-2 py-1 rounded border text-xs font-bold ${row.status.includes('Blocked') ? 'bg-[#ef4444]/10 text-[#ef4444] border-[#ef4444]/20' : 'bg-[#10b981]/10 text-[#10b981] border-[#10b981]/20'}`}>
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
