"use client"
import Navbar from '@/components/Navbar';

export default function FeaturesPage() {
  return (
    <div className="min-h-screen flex flex-col font-[family-name:var(--font-sans)] relative overflow-hidden bg-[#05050a] text-zinc-300">
      <div className="fixed inset-0 z-0 bg-[linear-gradient(to_right,#ffffff05_1px,transparent_1px),linear-gradient(to_bottom,#ffffff05_1px,transparent_1px)] bg-[size:4rem_4rem]"></div>
      
      <Navbar />

      <main className="flex flex-col items-center w-full z-10 pt-16 px-8 pb-32 max-w-6xl mx-auto">
        <div className="text-center mb-16">
          <h1 className="text-5xl font-black text-white mb-6 tracking-tight">Unprecedented Control.</h1>
          <p className="text-zinc-400 text-xl max-w-2xl mx-auto">
            Stop relying on prompt engineering for safety. Rely on mathematics, the Linux Kernel, and cryptographic object databases.
          </p>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 w-full">
          <div className="glass-card p-8 flex flex-col gap-4 hover:border-[var(--accent)] transition-all duration-300 group">
            <div className="text-4xl group-hover:scale-110 transition-transform origin-left">🛡️</div>
            <h3 className="text-xl font-bold text-white">eBPF Interceptors</h3>
            <p className="text-sm text-zinc-400 leading-relaxed">Agents hit mock environments instead of live databases. Our raw C-code eBPF programs hook into the kernel to silently drop DROP/DELETE SQL queries before they leave the container.</p>
          </div>
          
          <div className="glass-card p-8 flex flex-col gap-4 hover:border-[#8a2be2] transition-all duration-300 group">
            <div className="text-4xl group-hover:scale-110 transition-transform origin-left">⏪</div>
            <h3 className="text-xl font-bold text-white">Zero-Copy Time Travel</h3>
            <p className="text-sm text-zinc-400 leading-relaxed">StuntDouble takes a cryptographic snapshot of your working directory using Git's low-level core object database. If an AI ruins your code, `stuntdouble rewind` instantly restores it.</p>
          </div>
          
          <div className="glass-card p-8 flex flex-col gap-4 hover:border-[var(--accent)] transition-all duration-300 group">
            <div className="text-4xl group-hover:scale-110 transition-transform origin-left">🌐</div>
            <h3 className="text-xl font-bold text-white">StuntNet Swarms</h3>
            <p className="text-sm text-zinc-400 leading-relaxed">Orchestrate teams of autonomous agents inside a virtual intranet. They communicate over a custom Docker bridge network that hard-blocks all external internet access.</p>
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
            <div className="text-4xl group-hover:scale-110 transition-transform origin-left">☁️</div>
            <h3 className="text-xl font-bold text-white">Enterprise Cloud</h3>
            <p className="text-sm text-zinc-400 leading-relaxed">Sync telemetry, eBPF block events, and swarm activity directly to the Cloud dashboard. Keep your engineering team secure across the globe.</p>
          </div>
        </div>
      </main>

      <footer className="border-t border-zinc-800 w-full py-12 flex flex-col items-center justify-center text-zinc-500 text-sm z-10 bg-[#05050a] mt-auto">
        <div className="text-lg font-black text-white mb-2">Stunt<span className="text-[#00f0ff]">Double</span></div>
        <p>Built for the AI Agent Era. © 2026</p>
      </footer>
    </div>
  );
}
