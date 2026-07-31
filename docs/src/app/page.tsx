"use client"
import Navbar from '@/components/Navbar';

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-start font-[family-name:var(--font-sans)] relative overflow-hidden bg-[#05050a] text-zinc-300">
      
      {/* Background Effects */}
      <div className="fixed inset-0 z-0 bg-[linear-gradient(to_right,#ffffff05_1px,transparent_1px),linear-gradient(to_bottom,#ffffff05_1px,transparent_1px)] bg-[size:4rem_4rem]"></div>
      
      <Navbar />
      
      <main className="flex flex-col items-center w-full z-10 pt-20 px-8 pb-32">
        
        {/* Hero Section */}
        <section className="text-center w-full flex flex-col items-center gap-8 animate-float max-w-4xl">
          <div className="inline-flex px-4 py-1.5 rounded-full border border-[rgba(0,240,255,0.3)] bg-[rgba(0,240,255,0.05)] text-[var(--accent)] text-sm font-semibold tracking-widest uppercase shadow-[0_0_20px_rgba(0,240,255,0.15)]">
            v3.4.3 Engine Live • Zero-Copy Rewinds
          </div>
          <h1 className="text-6xl sm:text-8xl font-extrabold tracking-tight text-white leading-[1.1]">
            Deploy Agents. <br/>
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#00f0ff] to-[#8a2be2]">Without the Fear.</span>
          </h1>
          <p className="text-xl text-zinc-400 max-w-2xl mt-2 leading-relaxed font-light">
            StuntDouble is the ultimate Zero-Trust Isolation Engine. We wrap autonomous coding agents like Claude Code and Copilot in restricted containers with dropped capabilities, resource limits, zero-copy rewinds, and native macOS &amp; Windows kernel interceptors.
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

        {/* Detailed Command Reference */}
        <section className="w-full mt-32 max-w-4xl flex flex-col gap-8">
          <div className="text-center mb-8">
            <h2 className="text-4xl font-bold text-white mb-4">CLI Reference Manual</h2>
            <p className="text-zinc-400 text-lg">Master the StuntDouble ecosystem with our zero-trust orchestration commands.</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-[#111116] border border-zinc-800 p-8 rounded-3xl hover:border-[#00f0ff] transition-colors group">
              <h3 className="text-2xl font-bold text-white mb-2 font-mono group-hover:text-[#00f0ff] transition-colors">sd run &lt;agent&gt;</h3>
              <p className="text-zinc-400 leading-relaxed mb-4">Spawns an agent inside a capability-dropped container with scoped workspace mounts.</p>
              <ul className="text-sm text-zinc-500 space-y-2 font-mono">
                <li>--allow-unenforced-network : Acknowledge egress state</li>
                <li>--env (-e) : Inject custom runtime base image</li>
              </ul>
            </div>

            <div className="bg-[#111116] border border-zinc-800 p-8 rounded-3xl hover:border-[#8a2be2] transition-colors group">
              <h3 className="text-2xl font-bold text-white mb-2 font-mono group-hover:text-[#8a2be2] transition-colors">sd rewind</h3>
              <p className="text-zinc-400 leading-relaxed mb-4">Restores tracked workspace files to the pre-run snapshot using zero-copy git plumbing and removes untracked agent artifacts.</p>
            </div>

            <div className="bg-[#111116] border border-zinc-800 p-8 rounded-3xl hover:border-[#00f0ff] transition-colors group">
              <h3 className="text-2xl font-bold text-white mb-2 font-mono group-hover:text-[#00f0ff] transition-colors">sd swarm &lt;agents...&gt;</h3>
              <p className="text-zinc-400 leading-relaxed mb-4">Orchestrates multiple AI agents inside an isolated Docker bridge network (StuntNet) with synthetic internal service communication.</p>
            </div>

            <div className="bg-[#111116] border border-zinc-800 p-8 rounded-3xl hover:border-[#8a2be2] transition-colors group">
              <h3 className="text-2xl font-bold text-white mb-2 font-mono group-hover:text-[#8a2be2] transition-colors">sd record &lt;command&gt;</h3>
              <p className="text-zinc-400 leading-relaxed mb-4">Runs target commands under a Keploy sidecar to capture outbound database and API traffic into replayable mock definitions.</p>
            </div>
          </div>
        </section>
      </main>
      
      <footer className="border-t border-zinc-800 w-full py-12 flex flex-col items-center justify-center text-zinc-500 text-sm z-10 bg-[#05050a] mt-auto">
        <div className="text-lg font-black text-white mb-2">Stunt<span className="text-[#00f0ff]">Double</span></div>
        <p>Built for the AI Agent Era. © 2026</p>
      </footer>
    </div>
  );
}
