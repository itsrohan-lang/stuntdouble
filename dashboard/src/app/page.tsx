"use client"
import React, { useState, useEffect } from 'react';
import { Activity, ShieldAlert, Users, ServerCrash, Terminal, Lock, Globe } from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, AreaChart, Area } from 'recharts';

export default function Dashboard() {
  const [telemetry, setTelemetry] = useState({ total_runs: 0, blocked_commands: 0 });

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
          <ShieldAlert className="w-8 h-8 text-[#00f0ff]" />
          <span className="text-xl font-black text-white tracking-tight">StuntDouble <span className="text-zinc-500 font-medium">Control Plane</span></span>
        </div>
        <div className="flex items-center gap-6 text-sm font-medium">
          <a href="#" className="text-white">Overview</a>
          <a href="#" className="text-zinc-500 hover:text-zinc-300 transition">Policies</a>
          <a href="#" className="text-zinc-500 hover:text-zinc-300 transition">Audit Logs</a>
          <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-[#8a2be2] to-[#00f0ff] flex items-center justify-center text-white font-bold ml-4 shadow-[0_0_15px_rgba(0,240,255,0.3)]">
            CTO
          </div>
        </div>
      </nav>

      <main className="relative z-10 p-8 max-w-7xl mx-auto space-y-8">
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
              <h2 className="text-xl font-bold text-white">Blocked Outbound Connections</h2>
              <span className="px-3 py-1 bg-[#ef4444]/10 text-[#ef4444] text-xs font-bold rounded-full border border-[#ef4444]/20 uppercase tracking-wider">Live</span>
            </div>
            <div className="h-[300px] w-full">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={data}>
                  <defs>
                    <linearGradient id="colorBlocked" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#ef4444" stopOpacity={0.3}/>
                      <stop offset="95%" stopColor="#ef4444" stopOpacity={0}/>
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" stroke="#27272a" vertical={false} />
                  <XAxis dataKey="time" stroke="#52525b" tick={{fill: '#71717a'}} axisLine={false} tickLine={false} />
                  <YAxis stroke="#52525b" tick={{fill: '#71717a'}} axisLine={false} tickLine={false} />
                  <Tooltip 
                    contentStyle={{ backgroundColor: '#18181b', border: '1px solid #3f3f46', borderRadius: '12px', color: '#fff' }}
                    itemStyle={{ color: '#ef4444' }}
                  />
                  <Area type="monotone" dataKey="blocked" stroke="#ef4444" strokeWidth={3} fillOpacity={1} fill="url(#colorBlocked)" />
                </AreaChart>
              </ResponsiveContainer>
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
      </main>
    </div>
  );
}
