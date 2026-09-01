import type { ReactNode } from 'react'

import type { PanelStatus } from './types'

type PanelProps = { id: string; title: string; status?: PanelStatus; active?: boolean; children: ReactNode; className?: string }

const statusClass: Record<PanelStatus, string> = {
  idle: 'bg-white text-black',
  live: 'bg-black text-white',
  queued: 'bg-[repeating-linear-gradient(135deg,#fff_0,#fff_4px,#111_4px,#111_5px)] text-black',
  flushing: 'bg-black text-white',
}

export function Panel({ id, title, status = 'idle', active = false, children, className = '' }: PanelProps) {
  return <section className={`border-2 border-black bg-[#f8f8f2] shadow-[4px_4px_0_#111] transition ${active ? 'scale-[1.01] ring-4 ring-black animate-pulse' : ''} ${className}`}><header className="flex items-center justify-between border-b-2 border-black bg-white px-3 py-2 font-mono text-[10px] font-bold uppercase tracking-[0.18em]"><span>{id} // {title}</span><span className={`border border-black px-1.5 py-0.5 text-[9px] tracking-[0.12em] ${statusClass[status]}`}>{status}</span></header>{children}</section>
}

export function Metric({ label, value }: { label: string; value: string | number }) {
  return <div className="flex items-baseline justify-between gap-3 border-b border-dashed border-black/35 py-1.5 font-mono text-xs last:border-b-0"><span className="uppercase tracking-[0.12em] text-black/55">{label}</span><span className="font-bold text-black">{value}</span></div>
}
