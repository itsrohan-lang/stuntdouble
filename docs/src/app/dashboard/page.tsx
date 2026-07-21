"use client"
import { useEffect, useState } from 'react';
import Link from 'next/link';

interface TelemetryStats {
  total_runs: number;
  blocked_commands: number;
  last_run: string;
  status: string;
}

export default function Dashboard() {
  const [stats, setStats] = useState<TelemetryStats | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isPolling, setIsPolling] = useState(true);
  const [isClient, setIsClient] = useState(false);

  useEffect(() => {
    setIsClient(true);
    const fetchStats = async () => {
      try {
        const res = await fetch('http://localhost:8080/api/stats');
        if (!res.ok) throw new Error("API Offline");
        const data = await res.json();
        setStats(data);
        setError(null);
      } catch (err) {
        setError("Control Plane Offline. Run 'stuntdouble serve' locally.");
      }
    };

    fetchStats();
    const interval = setInterval(() => {
      if (isPolling) fetchStats();
    }, 2000);

    return () => clearInterval(interval);
  }, [isPolling]);

  return (
    <div className="min-h-screen flex flex-col font-[family-name:var(--font-sans)] bg-[#05050a] text-zinc-300 relative overflow-hidden">
      
      {/* Dynamic Background Effects */}
      <div className="fixed inset-0 z-0 bg-[linear-gradient(to_right,#ffffff05_1px,transparent_1px),linear-gradient(to_bottom,#ffffff05_1px,transparent_1px)] bg-[size:4rem_4rem]"></div>
      <div className={`fixed top-1/4 left-1/4 w-[500px] h-[500px] bg-[#00f0ff] opacity-[0.03] blur-[120px] rounded-full mix-blend-screen transition-all duration-1000 ${error ? 'bg-red-500 opacity-[0.05]' : 'bg-[#00f0ff]'}`}></div>

      {/* Navbar */}
      <nav className="w-full flex justify-between items-center py-4 px-8 z-20 border-b border-zinc-800/50 backdrop-blur-md bg-[#05050a]/80">
        <Link href="/">
          <div className="text-xl font-black tracking-tighter text-white cursor-pointer hover:opacity-80 transition">
            Stunt<span className="text-[#00f0ff]">Double</span> <span className="text-zinc-500 font-medium text-sm ml-2">Control Plane</span>
          </div>
        </Link>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <span className="text-xs font-semibold text-zinc-400 tracking-widest uppercase">Status</span>
            {error ? (
              <div className="flex items-center gap-2 bg-red-950/30 border border-red-900/50 px-3 py-1 rounded-full">
                <div className="w-2 h-2 rounded-full bg-red-500 animate-pulse"></div>
                <span className="text-xs font-bold text-red-400">Offline</span>
              </div>
            ) : (
              <div className="flex items-center gap-2 bg-[#00f0ff]/10 border border-[#00f0ff]/20 px-3 py-1 rounded-full">
                <div className="w-2 h-2 rounded-full bg-[#00f0ff] animate-pulse" style={{boxShadow: "0 0 10px #00f0ff"}}></div>
                <span className="text-xs font-bold text-[#00f0ff]">Active</span>
              </div>
            )}
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="flex-1 flex flex-col items-center justify-start p-8 z-10 w-full max-w-6xl mx-auto pt-16">
        
        {/* Header */}
        <div className="w-full mb-12 flex justify-between items-end">
          <div>
            <h1 className="text-4xl font-black text-white mb-2 tracking-tight">Mission Control</h1>
            <p className="text-zinc-400">Live telemetry from the eBPF kernel layer & Docker sandbox.</p>
          </div>
          <button 
            onClick={() => setIsPolling(!isPolling)}
            className={`px-4 py-2 rounded-md text-xs font-bold border transition-all ${isPolling ? 'border-zinc-700 text-zinc-400 hover:bg-zinc-900' : 'border-[#00f0ff]/50 text-[#00f0ff] bg-[#00f0ff]/10 hover:bg-[#00f0ff]/20'}`}
          >
            {isPolling ? 'Pause Telemetry' : 'Resume Telemetry'}
          </button>
        </div>

        {/* Dashboard Grid */}
        <div className="w-full grid grid-cols-1 md:grid-cols-3 gap-6">
          
          {/* Total Runs Card */}
          <div className="glass-card p-6 flex flex-col relative overflow-hidden group">
            <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-[#8a2be2] to-transparent opacity-50"></div>
            <h3 className="text-sm font-semibold text-zinc-400 tracking-wider uppercase mb-1">Total Agent Executions</h3>
            <div className="text-5xl font-black text-white my-4">
              {stats?.total_runs !== undefined ? stats.total_runs : '-'}
            </div>
            <p className="text-xs text-zinc-500 mt-auto">Isolated runs inside Docker MicroVMs.</p>
          </div>

          {/* Blocked Commands Card */}
          <div className="glass-card p-6 flex flex-col relative overflow-hidden group border-red-900/30">
            <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-red-500 to-transparent opacity-50"></div>
            <div className="absolute -right-10 -bottom-10 w-40 h-40 bg-red-500/10 rounded-full blur-2xl group-hover:bg-red-500/20 transition-all"></div>
            <h3 className="text-sm font-semibold text-zinc-400 tracking-wider uppercase mb-1">eBPF Packet Drops</h3>
            <div className="text-5xl font-black text-white my-4 text-red-400 drop-shadow-[0_0_15px_rgba(248,113,113,0.3)]">
              {stats?.blocked_commands !== undefined ? stats.blocked_commands : '-'}
            </div>
            <p className="text-xs text-zinc-500 mt-auto text-red-200/50">Malicious queries intercepted at kernel level.</p>
          </div>

          {/* System Status Card */}
          <div className="glass-card p-6 flex flex-col relative overflow-hidden">
            <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-[#00f0ff] to-transparent opacity-50"></div>
            <h3 className="text-sm font-semibold text-zinc-400 tracking-wider uppercase mb-1">Kernel Status</h3>
            <div className="flex items-center gap-3 my-4">
              <div className="text-2xl font-bold text-white">
                {error ? "Disconnected" : stats?.status || "Analyzing..."}
              </div>
            </div>
            {stats?.last_run && (
              <p className="text-xs text-zinc-500 mt-auto">
                Last activity: {new Date(stats.last_run).toLocaleTimeString()}
              </p>
            )}
          </div>
        </div>

        {/* Console Output Mock */}
        <div className="w-full mt-10 glass-card p-0 overflow-hidden border-zinc-800 flex flex-col">
          <div className="bg-[#0a0a0f] border-b border-zinc-800 px-4 py-3 flex items-center gap-2">
            <div className="w-3 h-3 rounded-full bg-red-500/50"></div>
            <div className="w-3 h-3 rounded-full bg-yellow-500/50"></div>
            <div className="w-3 h-3 rounded-full bg-green-500/50"></div>
            <span className="text-xs font-mono text-zinc-500 ml-2">/var/log/stuntdouble/warden.log</span>
          </div>
          <div className="p-6 font-mono text-sm h-64 overflow-y-auto flex flex-col gap-2">
            {isClient && (
              <>
                <div className="text-zinc-500">[{new Date().toISOString()}] Warden initialized. Awaiting API connection...</div>
                {error && <div className="text-red-400">[{new Date().toISOString()}] ERROR: Connection to localhost:8080 refused. Start 'stuntdouble serve'.</div>}
                {!error && stats && (
                  <>
                    <div className="text-green-400">[{new Date().toISOString()}] Connected to Control Plane API on port 8080.</div>
                    <div className="text-[#00f0ff]">[{new Date().toISOString()}] Listening for eBPF ring buffer events...</div>
                  </>
                )}
              </>
            )}
          </div>
        </div>

      </main>
    </div>
  );
}
