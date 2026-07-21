"use client"
import Head from 'next/head';

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-start font-[family-name:var(--font-sans)] relative overflow-hidden bg-[#05050a] text-zinc-300">
      
      {/* Background Effects */}
      <div className="fixed inset-0 z-0 bg-[linear-gradient(to_right,#ffffff05_1px,transparent_1px),linear-gradient(to_bottom,#ffffff05_1px,transparent_1px)] bg-[size:4rem_4rem]"></div>
      
      <nav className="w-full max-w-6xl flex justify-between items-center py-6 px-8 z-20 border-b border-zinc-800/50 backdrop-blur-md sticky top-0 bg-[#05050a]/80">
        <div className="text-2xl font-black tracking-tighter text-white">Stunt<span className="text-[#00f0ff]">Double</span></div>
        <div className="flex gap-6 items-center">
          <a href="#features" className="text-sm font-medium hover:text-white transition hidden sm:block">Features</a>
          <a href="#commands" className="text-sm font-medium hover:text-white transition hidden sm:block">Commands</a>
          <a className="bg-white text-black px-5 py-2 rounded-full text-sm font-bold hover:bg-zinc-200 transition shadow-[0_0_15px_rgba(255,255,255,0.3)]" href="https://github.com/itsrohan-lang/stuntdouble" target="_blank" rel="noopener noreferrer">
            GitHub ↗
          </a>
        </div>
      </nav>

      <main className="flex flex-col items-center w-full z-10 pt-20 px-8 pb-32">
        
        {/* Hero Section */}
        <section className="text-center w-full flex flex-col items-center gap-8 animate-float max-w-4xl">
          <div className="inline-flex px-4 py-1.5 rounded-full border border-[rgba(0,240,255,0.3)] bg-[rgba(0,240,255,0.05)] text-[var(--accent)] text-sm font-semibold tracking-widest uppercase shadow-[0_0_20px_rgba(0,240,255,0.15)]">
            v2.0 Engine Live • Native eBPF Hooks
          </div>
          <h1 className="text-6xl sm:text-8xl font-extrabold tracking-tight text-white leading-[1.1]">
            Deploy Agents. <br/>
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#00f0ff] to-[#8a2be2]">Without the Fear.</span>
          </h1>
          <p className="text-xl text-zinc-400 max-w-2xl mt-2 leading-relaxed font-light">
            StuntDouble is the ultimate Zero-Trust Isolation Engine. We wrap autonomous tools like Claude Code and Copilot in hardware-level MicroVMs and drop destructive network calls using raw kernel eBPF hooks.
          </p>
          
          <div className="flex gap-4 mt-4">
            <a className="btn-primary rounded-full px-8 py-4 font-bold text-md tracking-wide cursor-pointer" href="#quickstart">
              Install Native CLI
            </a>
            <a className="rounded-full px-8 py-4 font-bold text-md tracking-wide text-zinc-300 border border-zinc-700 hover:bg-zinc-800 transition-colors cursor-pointer flex items-center gap-2" href="https://github.com/itsrohan-lang/stuntdouble" target="_blank" rel="noopener noreferrer">
              Read the Docs
            </a>
          </div>
        </section>

        {/* Installation Snippet */}
        <section id="quickstart" className="w-full mt-24 max-w-3xl flex flex-col items-center">
          <div className="glass-card p-2 w-full flex items-center justify-between border-zinc-800 bg-[#0d0d12]/80 backdrop-blur-xl pl-6 rounded-2xl shadow-2xl">
            <code className="text-zinc-300 font-mono text-sm overflow-x-auto whitespace-nowrap hide-scrollbar">
              <span className="text-[#00f0ff]">curl</span> -sSL https://raw.githubusercontent.com/itsrohan-lang/stuntdouble/main/install.sh | bash
            </code>
            <button className="bg-zinc-800 hover:bg-zinc-700 text-white px-5 py-2.5 rounded-xl text-sm font-medium transition cursor-pointer ml-4 flex-shrink-0" onClick={() => navigator.clipboard.writeText('curl -sSL https://raw.githubusercontent.com/itsrohan-lang/stuntdouble/main/install.sh | bash')}>
              Copy
            </button>
          </div>
        </section>

        {/* Feature Grid */}
        <section id="features" className="w-full max-w-6xl mt-32">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-white mb-4">Unprecedented Control.</h2>
            <p className="text-zinc-400 text-lg">Stop relying on prompt engineering for safety. Rely on mathematics and the Linux Kernel.</p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div className="glass-card p-8 flex flex-col gap-4 hover:border-[var(--accent)] transition-all duration-300 group">
              <div className="text-4xl group-hover:scale-110 transition-transform origin-left">🛡️</div>
              <h3 className="text-xl font-bold text-white">eBPF Interceptors</h3>
              <p className="text-sm text-zinc-400 leading-relaxed">Agents hit mock environments instead of live databases. Our raw C-code eBPF programs hook into the kernel to silently drop DROP/DELETE SQL queries before they leave the container.</p>
            </div>
            
            <div className="glass-card p-8 flex flex-col gap-4 hover:border-[#8a2be2] transition-all duration-300 group">
              <div className="text-4xl group-hover:scale-110 transition-transform origin-left">⏪</div>
              <h3 className="text-xl font-bold text-white">Time-Travel Rewind</h3>
              <p className="text-sm text-zinc-400 leading-relaxed">If an autonomous agent hallucinating deletes your workspace, run `stuntdouble rewind` to instantly restore the ZFS filesystem snapshot to exactly 5 minutes ago.</p>
            </div>
            
            <div className="glass-card p-8 flex flex-col gap-4 hover:border-[var(--accent)] transition-all duration-300 group">
              <div className="text-4xl group-hover:scale-110 transition-transform origin-left">🌐</div>
              <h3 className="text-xl font-bold text-white">StuntNet Swarms</h3>
              <p className="text-sm text-zinc-400 leading-relaxed">Orchestrate teams of QA and Developer agents inside a virtual intranet. They can communicate over mocked sockets, completely isolated from the real web.</p>
            </div>

            <div className="glass-card p-8 flex flex-col gap-4 hover:border-[#8a2be2] transition-all duration-300 group">
              <div className="text-4xl group-hover:scale-110 transition-transform origin-left">⚡</div>
              <h3 className="text-xl font-bold text-white">Native Docker SDK</h3>
              <p className="text-sm text-zinc-400 leading-relaxed">No fragile bash scripts. StuntDouble talks directly to the Docker Daemon Unix socket in Go to stream TTY inputs interactively and drop all root capabilities instantly.</p>
            </div>

            <div className="glass-card p-8 flex flex-col gap-4 hover:border-[var(--accent)] transition-all duration-300 group">
              <div className="text-4xl group-hover:scale-110 transition-transform origin-left">🔐</div>
              <h3 className="text-xl font-bold text-white">STP Governance</h3>
              <p className="text-sm text-zinc-400 leading-relaxed">The Universal Stunt Protocol. Foundational LLMs query our local HTTP server to verify cryptographic sandbox attestations before executing code.</p>
            </div>

            <div className="glass-card p-8 flex flex-col gap-4 hover:border-[#8a2be2] transition-all duration-300 group">
              <div className="text-4xl group-hover:scale-110 transition-transform origin-left">👁️</div>
              <h3 className="text-xl font-bold text-white">The Warden</h3>
              <p className="text-sm text-zinc-400 leading-relaxed">An autonomous defensive AI runs parallel to your agent, monitoring syscall telemetry. If an escape is detected, it generates zero-day eBPF patches on the fly.</p>
            </div>
          </div>
        </section>

        {/* CLI Reference */}
        <section id="commands" className="w-full max-w-5xl mt-40">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-white mb-4">Command Line Interface</h2>
            <p className="text-zinc-400 text-lg">A powerful, intuitive toolkit for managing agent safety.</p>
          </div>
          
          <div className="flex flex-col gap-4">
            {[
              { cmd: 'stuntdouble run claude', desc: 'Wraps Claude Code in a hardware-isolated microVM with native TTY streaming.' },
              { cmd: 'stuntdouble warden', desc: 'Starts the autonomous defensive AI monitor to oversee agent actions.' },
              { cmd: 'stuntdouble rewind 5', desc: 'Restores the workspace to the exact state it was in 5 minutes ago.' },
              { cmd: 'stuntdouble protocol start', desc: 'Boots the local STP attestation server on port 4438 for LLM handshakes.' },
              { cmd: 'stuntdouble stats', desc: 'Displays real-time eBPF network drops and packet isolation telemetry.' },
            ].map((item, i) => (
              <div key={i} className="glass-card p-6 flex flex-col sm:flex-row items-start sm:items-center justify-between border-zinc-800/50 hover:bg-zinc-900/50 transition">
                <code className="text-[#00f0ff] font-mono text-lg font-semibold">{item.cmd}</code>
                <p className="text-zinc-400 text-sm mt-2 sm:mt-0 sm:text-right max-w-sm">{item.desc}</p>
              </div>
            ))}
          </div>
        </section>

      </main>
      
      <footer className="border-t border-zinc-800 w-full py-12 flex flex-col items-center justify-center text-zinc-500 text-sm z-10 bg-[#05050a]">
        <div className="text-lg font-black text-white mb-2">Stunt<span className="text-[#00f0ff]">Double</span></div>
        <p>Built for the AI Agent Era. © 2026</p>
      </footer>
    </div>
  );
}
