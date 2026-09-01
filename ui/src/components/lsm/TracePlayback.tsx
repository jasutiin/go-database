import { ChevronRight, X } from 'lucide-react'

import type { TraceEvent } from './types'

type TracePlaybackProps = {
  current: TraceEvent
  remaining: number
  onNext: () => void
  onFinish: () => void
}

export function TracePlayback({ current, remaining, onNext, onFinish }: TracePlaybackProps) {
  const isFinalEvent = remaining === 0

  return (
    <article className="w-full max-w-xl border-2 border-black bg-[#f8f8f2] shadow-[7px_7px_0_#111]">
      <header className="flex items-center justify-between border-b-2 border-black bg-black px-4 py-3 text-white">
        <p className="font-mono text-[10px] font-bold uppercase tracking-[0.18em]">trace card // {current.target}</p>
        <span className="border border-white px-2 py-1 font-mono text-[10px] uppercase">{remaining + 1} step{remaining === 0 ? '' : 's'} left</span>
      </header>
      <div className="p-5">
        <p className="font-mono text-[10px] font-bold uppercase tracking-[0.15em] text-black/55">{current.source} event</p>
        <h2 className="mt-3 font-mono text-2xl font-black uppercase tracking-[-0.05em]">{current.title}</h2>
        <p className="mt-4 border-l-2 border-black pl-3 font-mono text-sm leading-6">{current.detail}</p>
      </div>
      <footer className="flex justify-end gap-2 border-t-2 border-black bg-white p-3">
        {isFinalEvent ? <button type="button" onClick={onFinish} className="flex items-center gap-2 border-2 border-black bg-black px-3 py-2 font-mono text-xs font-bold uppercase text-white hover:bg-white hover:text-black"><X size={14} /> finish trace</button> : <button type="button" onClick={onNext} className="flex items-center gap-2 border-2 border-black bg-black px-3 py-2 font-mono text-xs font-bold uppercase text-white hover:bg-white hover:text-black">next event <ChevronRight size={14} /></button>}
      </footer>
    </article>
  )
}
