import Link from 'next/link';

export default function Navbar() {
  return (
    <nav className="w-full max-w-6xl flex justify-between items-center py-6 px-8 z-20 border-b border-zinc-800/50 backdrop-blur-md sticky top-0 bg-[#05050a]/80 mx-auto">
      <Link href="/">
        <div className="text-2xl font-black tracking-tighter text-white cursor-pointer hover:opacity-80 transition">
          Stunt<span className="text-[#00f0ff]">Double</span>
        </div>
      </Link>
      <div className="flex gap-6 items-center">
        <Link href="/features" className="text-sm font-medium hover:text-white transition hidden sm:block text-zinc-300">
          Features
        </Link>
        <Link href="/commands" className="text-sm font-medium hover:text-white transition hidden sm:block text-zinc-300">
          Commands
        </Link>
        <Link href="/dashboard" className="text-sm font-bold text-[#00f0ff] hover:text-white transition hidden sm:block">
          Dashboard
        </Link>
        <a 
          className="bg-white text-black px-5 py-2 rounded-full text-sm font-bold hover:bg-zinc-200 transition shadow-[0_0_15px_rgba(255,255,255,0.3)]" 
          href="https://github.com/itsrohan-lang/stuntdouble" 
          target="_blank" 
          rel="noopener noreferrer"
        >
          GitHub ↗
        </a>
      </div>
    </nav>
  );
}
