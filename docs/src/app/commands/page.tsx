"use client"
import Navbar from '@/components/Navbar';

export default function CommandsPage() {
  return (
    <div className="min-h-screen flex flex-col font-[family-name:var(--font-sans)] relative overflow-hidden bg-[#05050a] text-zinc-300">
      <div className="fixed inset-0 z-0 bg-[linear-gradient(to_right,#ffffff05_1px,transparent_1px),linear-gradient(to_bottom,#ffffff05_1px,transparent_1px)] bg-[size:4rem_4rem]"></div>
      
      <Navbar />

      <main className="flex flex-col items-center w-full z-10 pt-16 px-8 pb-32 max-w-5xl mx-auto">
        <div className="text-center mb-16">
          <h1 className="text-5xl font-black text-white mb-6 tracking-tight">Command Line Interface</h1>
          <p className="text-zinc-400 text-xl max-w-2xl mx-auto">
            A powerful, intuitive toolkit for managing agent safety, swarms, and time-travel.
          </p>
        </div>
        
        <div className="flex flex-col gap-6 w-full">
          {[
            { cmd: 'stuntdouble run [agent]', desc: 'Wraps an agent (like `claude` or `aider`) in a hardware-isolated microVM with native TTY streaming. Takes a zero-copy snapshot of your files before executing.' },
            { cmd: 'stuntdouble rewind', desc: 'Instantly restores your entire workspace to the exact state it was in before the last agent ran. Cleans up all rogue files.' },
            { cmd: 'stuntdouble swarm [agents...]', desc: 'Orchestrates a team of autonomous agents inside StuntNet—an internal, internet-blocked Docker bridge network.' },
            { cmd: 'stuntdouble serve', desc: 'Boots the local Control Plane API on port 8080 to power the visual web dashboard.' },
            { cmd: 'stuntdouble login [token]', desc: 'Authenticates with StuntDouble Enterprise Cloud to sync global telemetry for your engineering team.' },
          ].map((item, i) => (
            <div key={i} className="glass-card p-8 flex flex-col lg:flex-row items-start lg:items-center justify-between border-zinc-800/50 hover:bg-zinc-900/50 transition">
              <code className="text-[#00f0ff] font-mono text-xl font-bold bg-[#00f0ff]/10 px-4 py-2 rounded-lg mb-4 lg:mb-0">{item.cmd}</code>
              <p className="text-zinc-400 text-base lg:text-right max-w-xl leading-relaxed">{item.desc}</p>
            </div>
          ))}
        </div>
      </main>

      <footer className="border-t border-zinc-800 w-full py-12 flex flex-col items-center justify-center text-zinc-500 text-sm z-10 bg-[#05050a] mt-auto">
        <div className="text-lg font-black text-white mb-2">Stunt<span className="text-[#00f0ff]">Double</span></div>
        <p>Built for the AI Agent Era. © 2026</p>
      </footer>
    </div>
  );
}
