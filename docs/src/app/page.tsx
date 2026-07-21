import Head from 'next/head';

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-8 sm:p-20 font-[family-name:var(--font-sans)] relative overflow-hidden">
      
      {/* Background Grid Pattern */}
      <div className="absolute inset-0 z-0 bg-[linear-gradient(to_right,#ffffff05_1px,transparent_1px),linear-gradient(to_bottom,#ffffff05_1px,transparent_1px)] bg-[size:4rem_4rem]"></div>

      <main className="flex flex-col gap-12 items-center sm:items-start max-w-5xl w-full z-10 pt-10">
        
        {/* Hero Section */}
        <section className="text-center sm:text-left w-full flex flex-col items-center sm:items-start gap-6 animate-float mt-10">
          <div className="inline-block px-4 py-1.5 rounded-full border border-[rgba(0,240,255,0.3)] bg-[rgba(0,240,255,0.05)] text-[var(--accent)] text-sm font-semibold tracking-widest uppercase mb-4 shadow-[0_0_15px_rgba(0,240,255,0.2)]">
            Universal AI Governance Protocol
          </div>
          <h1 className="text-5xl sm:text-7xl font-extrabold tracking-tight text-white leading-tight">
            Deploy Agents. <br/>
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#00f0ff] to-[#8a2be2]">Without the Fear.</span>
          </h1>
          <p className="text-lg text-zinc-400 max-w-2xl mt-4 leading-relaxed">
            StuntDouble is the ultimate Zero-Trust Isolation Engine for AI coding agents. We wrap tools like Claude Code and GitHub Copilot in hardware-level MicroVMs and eBPF network mocks.
          </p>
          
          <div className="flex gap-4 mt-6">
            <a className="btn-primary rounded-full px-8 py-3 font-semibold text-sm tracking-wide cursor-pointer" href="#quickstart">
              Get Started
            </a>
            <a className="rounded-full px-8 py-3 font-semibold text-sm tracking-wide text-zinc-300 border border-zinc-700 hover:bg-zinc-800 transition-colors cursor-pointer" href="https://github.com/itsrohan-lang/stuntdouble" target="_blank" rel="noopener noreferrer">
              View on GitHub
            </a>
          </div>
        </section>

        {/* Feature Cards */}
        <section className="grid grid-cols-1 sm:grid-cols-3 gap-6 w-full mt-16">
          <div className="glass-card p-8 flex flex-col gap-4 hover:border-[var(--accent)] transition-colors duration-300">
            <div className="text-3xl">🛡️</div>
            <h3 className="text-xl font-bold text-white">eBPF Interceptors</h3>
            <p className="text-sm text-zinc-400">Agents hit mock environments instead of live production databases, completely invisibly.</p>
          </div>
          <div className="glass-card p-8 flex flex-col gap-4 hover:border-[#8a2be2] transition-colors duration-300">
            <div className="text-3xl">⏪</div>
            <h3 className="text-xl font-bold text-white">Time-Travel Rewind</h3>
            <p className="text-sm text-zinc-400">If an agent deletes your workspace, instantly rewind the ZFS snapshot to 5 minutes ago.</p>
          </div>
          <div className="glass-card p-8 flex flex-col gap-4 hover:border-[var(--accent)] transition-colors duration-300">
            <div className="text-3xl">🌐</div>
            <h3 className="text-xl font-bold text-white">StuntNet Swarms</h3>
            <p className="text-sm text-zinc-400">Orchestrate teams of QA and Dev agents inside a virtual intranet isolated from the real web.</p>
          </div>
        </section>

        {/* Code Snippet */}
        <section id="quickstart" className="w-full mt-20">
          <h2 className="text-3xl font-bold text-white mb-6">Installation</h2>
          <div className="glass-card p-6 w-full max-w-3xl overflow-x-auto border-zinc-800 bg-[#0d0d12]">
            <pre className="text-zinc-300 font-mono text-sm leading-loose">
              <code className="text-zinc-500"># 1. Install the CLI globally</code><br/>
              <span className="text-[#00f0ff]">npx</span> stuntdouble-sandbox-cli init<br/><br/>
              <code className="text-zinc-500"># 2. Run your AI Agent securely</code><br/>
              <span className="text-[#00f0ff]">stuntdouble</span> run claude<br/><br/>
              <code className="text-zinc-500"># 3. Check safety telemetry logs</code><br/>
              <span className="text-[#00f0ff]">stuntdouble</span> stats
            </pre>
          </div>
        </section>

      </main>
      
      <footer className="mt-32 border-t border-zinc-800 w-full py-8 text-center text-zinc-500 text-sm z-10">
        StuntDouble © 2026. Built for the AI Agent Era.
      </footer>
    </div>
  );
}
