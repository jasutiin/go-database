import { ArrowRight, Database, HardDrive, Layers3, ListTree, Radio, Workflow } from 'lucide-react'
import { useCallback, useRef, useState } from 'react'
import type { ReactNode } from 'react'

import { CommandDeck } from './CommandDeck'
import { mockLsmSnapshot } from './mock-data'
import { Metric, Panel } from './Panel'
import { applyTraceEvent, createCompactionTrace, createOperationTrace } from './simulation'
import { TracePlayback } from './TracePlayback'
import type { CompactionSnapshot, ImmutableMemtableSnapshot, LsmDashboardSnapshot, MemtableSnapshot, SkipListSnapshot, SstableLevelSnapshot, StorageOperation, TraceEvent, WalSnapshot } from './types'

export function LsmDashboard() {
  const [snapshot, setSnapshot] = useState<LsmDashboardSnapshot>(mockLsmSnapshot)
  const [queue, setQueue] = useState<TraceEvent[]>([])
  const [current, setCurrent] = useState<TraceEvent | null>(null)
  const [commandOpen, setCommandOpen] = useState(true)
  const [traceOpen, setTraceOpen] = useState(false)
  const sequence = useRef(10492)

  const startTrace = useCallback((events: TraceEvent[]) => {
    const [first, ...remaining] = events
    if (!first) return
    setSnapshot((state) => applyTraceEvent(state, first))
    setCurrent(first)
    setQueue(remaining)
    setCommandOpen(false)
    setTraceOpen(true)
  }, [])

  const enqueueOperation = useCallback((operation: StorageOperation) => {
    sequence.current += 1
    startTrace(createOperationTrace(operation, sequence.current))
  }, [startTrace])

  const enqueueCompaction = useCallback(() => {
    sequence.current += 1
    const events = createCompactionTrace(sequence.current)
    if (current) {
      setQueue((queued) => [...queued, ...events])
      return
    }
    startTrace(events)
  }, [current, startTrace])

  const nextTraceEvent = useCallback(() => {
    setQueue((events) => {
      const [next, ...remaining] = events
      if (!next) return events
      setSnapshot((state) => applyTraceEvent(state, next))
      setCurrent(next)
      return remaining
    })
  }, [])

  function finishTrace() {
    setCurrent(null)
    setTraceOpen(false)
    setCommandOpen(true)
  }

  const activeTarget = current?.target

  return <main className="min-h-screen bg-[#e9e9e1] px-4 py-5 text-black sm:px-6 lg:px-8"><div className="mx-auto max-w-[1500px]"><DashboardHeader queueLength={queue.length} /><section className="mt-5 grid gap-5 xl:grid-cols-[0.85fr_1.65fr_0.85fr]"><WalPanel wal={snapshot.wal} active={activeTarget === 'wal'} /><ActiveMemtablePanel memtable={snapshot.activeMemtable} skipList={snapshot.skipList} active={activeTarget === 'memtable' || activeTarget === 'skip-list'} /><ImmutableMemtablePanel memtables={snapshot.immutableMemtables} active={activeTarget === 'immutable'} /></section><div className="my-5 flex items-center gap-2 overflow-hidden font-mono text-[10px] font-bold uppercase tracking-[0.16em] text-black/60"><span className="shrink-0">write path</span><span className="h-px w-full bg-black/40" /><ArrowRight className="shrink-0" size={16} strokeWidth={2.5} /><span className="shrink-0">persistent levels</span></div><section className="grid gap-5 lg:grid-cols-[0.9fr_1.6fr]"><CompactionPanel compaction={snapshot.compaction} active={activeTarget === 'compaction'} /><SstablesPanel levels={snapshot.sstableLevels} active={activeTarget === 'sstables'} /></section></div>{commandOpen && <Modal><CommandDeck onOperation={enqueueOperation} onCompaction={enqueueCompaction} /></Modal>}{traceOpen && current && <Modal><TracePlayback current={current} remaining={queue.length} onNext={nextTraceEvent} onFinish={finishTrace} /></Modal>}</main>
}

function DashboardHeader({ queueLength }: { queueLength: number }) {
  return <header className="border-2 border-black bg-black px-4 py-4 text-[#f8f8f2] shadow-[5px_5px_0_#111] sm:px-5"><div className="flex flex-col justify-between gap-4 lg:flex-row lg:items-center"><div className="flex items-center gap-3"><div className="grid h-10 w-10 place-items-center border-2 border-[#f8f8f2]"><Database size={22} strokeWidth={2.5} /></div><div><p className="font-mono text-[10px] font-bold uppercase tracking-[0.24em] text-white/65">db-engine // storage monitor</p><h1 className="font-mono text-2xl font-black uppercase tracking-[-0.06em] sm:text-3xl">lsm tree</h1></div></div><div className="flex flex-wrap items-center gap-2 font-mono text-[10px] font-bold uppercase tracking-[0.12em]"><span className="flex items-center gap-2 border border-white px-2 py-1.5"><Radio size={13} fill="currentColor" /> trace mode</span><span className="border border-white px-2 py-1.5">{queueLength} waiting</span></div></div></header>
}

function Modal({ children }: { children: ReactNode }) {
  return <div className="fixed inset-0 z-50 grid place-items-center bg-black/75 p-4"><div className="w-full max-w-5xl">{children}</div></div>
}

function WalPanel({ wal, active }: { wal: WalSnapshot; active: boolean }) { return <Panel id="01" title="write ahead log" status={wal.status} active={active}><div className="p-3"><p className="mb-3 font-mono text-sm font-black uppercase">{wal.segment}</p><div className="mb-4 h-3 border border-black bg-white p-[2px]"><div className="h-full w-[18%] bg-black" /></div><Metric label="records" value={wal.entries} /><Metric label="used" value={`${wal.usedBytes} / ${wal.capacity}`} /></div><div className="border-t-2 border-black bg-white p-3"><p className="mb-2 font-mono text-[10px] font-bold uppercase tracking-[0.16em]">latest records</p><ol className="space-y-1 font-mono text-[11px]">{wal.records.map((record, index) => <li key={`${record}-${index}`} className="flex gap-2"><span className="text-black/45">0{index + 1}</span><span>{record}</span></li>)}</ol></div></Panel> }

function ActiveMemtablePanel({ memtable, skipList, active }: { memtable: MemtableSnapshot; skipList: SkipListSnapshot; active: boolean }) { return <Panel id="02" title="active memtable" status={memtable.status} active={active}><div className="grid gap-4 p-3 lg:grid-cols-[190px_1fr]"><div><div className="mb-3 flex items-center gap-2 font-mono text-sm font-black uppercase"><Layers3 size={17} /> mutable set</div><Metric label="entries" value={memtable.entries} /><Metric label="memory" value={`${memtable.usedBytes} / ${memtable.limit}`} /><Metric label="rate" value={memtable.writeRate} /><div className="mt-4 border-2 border-black bg-white p-2 font-mono text-[10px] uppercase leading-relaxed"><span className="font-bold">index:</span> probabilistic skip list<br /><span className="font-bold">flush:</span> at 64 kb</div></div><SkipListPanel skipList={skipList} /></div></Panel> }

function SkipListPanel({ skipList }: { skipList: SkipListSnapshot }) { return <div className="border-2 border-black bg-white"><div className="flex items-center justify-between border-b-2 border-black px-3 py-2 font-mono text-[10px] font-bold uppercase tracking-[0.15em]"><span className="flex items-center gap-2"><ListTree size={14} /> skip list</span><span>{skipList.lookupSteps} lookup steps</span></div><div className="space-y-3 overflow-x-auto p-3">{skipList.levels.map((nodes, levelIndex) => <div key={`level-${levelIndex}`} className="flex min-w-max items-center gap-1.5 font-mono text-[10px]"><span className="w-12 border border-black bg-black px-1 py-1 text-center text-white">l{skipList.levels.length - levelIndex}</span><span className="text-lg leading-none">→</span>{nodes.map((node, index) => <div key={`${levelIndex}-${node}`} className="flex items-center gap-1.5"><span className={`border border-black px-2 py-1 ${node.startsWith('†') ? 'bg-black text-white' : 'bg-[#f8f8f2]'}`}>{node}</span>{index < nodes.length - 1 && <span className="text-lg leading-none">→</span>}</div>)}</div>)}</div><div className="border-t-2 border-black px-3 py-2 font-mono text-[10px] uppercase tracking-[0.12em] text-black/65">inverse node = tombstone</div></div> }

function ImmutableMemtablePanel({ memtables, active }: { memtables: ImmutableMemtableSnapshot[]; active: boolean }) { return <Panel id="03" title="immutable queue" status="queued" active={active}><div className="p-3"><p className="mb-3 font-mono text-sm font-black uppercase">flush queue</p><div className="space-y-3">{memtables.map((memtable, index) => <article key={memtable.id} className="border-2 border-black bg-white p-3"><div className="mb-2 flex items-center justify-between font-mono text-[10px] font-bold uppercase tracking-[0.12em]"><span>0{index + 1} // {memtable.id}</span><span className={memtable.status === 'flushing' ? 'bg-black px-1.5 py-0.5 text-white' : 'border border-black px-1.5 py-0.5'}>{memtable.status}</span></div><div className="flex justify-between font-mono text-xs"><span>{memtable.entries} entries</span><span>{memtable.size}</span></div></article>)}</div></div></Panel> }

function CompactionPanel({ compaction, active }: { compaction: CompactionSnapshot; active: boolean }) { return <Panel id="04" title="compaction" status={compaction.status} active={active}><div className="p-4"><div className="mb-5 flex items-center gap-2 font-mono text-sm font-black uppercase"><Workflow size={18} /> merge worker</div><div className="grid grid-cols-[1fr_auto_1fr] items-center gap-3 border-2 border-black bg-white p-3 font-mono text-[10px] uppercase"><div className="space-y-2"><div className="border border-black bg-black px-2 py-2 text-white">l0 / source a</div><div className="border border-black px-2 py-2">l0 / source b</div></div><ArrowRight size={22} strokeWidth={2.5} /><div className="border-2 border-black px-2 py-5 text-center font-bold">output<br />l1</div></div><div className="mt-4 grid grid-cols-2 gap-x-5"><Metric label="debt" value={compaction.debt} /><Metric label="last run" value={compaction.lastRun} /><Metric label="inputs" value={compaction.inputTables} /><Metric label="outputs" value={compaction.outputTables} /></div></div></Panel> }

function SstablesPanel({ levels, active }: { levels: SstableLevelSnapshot[]; active: boolean }) { return <Panel id="05" title="sorted string tables" status="idle" active={active}><div className="p-4"><div className="mb-4 flex items-center gap-2 font-mono text-sm font-black uppercase"><HardDrive size={18} /> persistent levels</div><div className="space-y-4">{levels.map((level) => <div key={level.level} className="grid grid-cols-[38px_1fr] gap-3"><div className="border-2 border-black bg-black py-2 text-center font-mono text-xs font-black uppercase text-white">{level.level}</div><div className="flex flex-wrap gap-2">{level.tables.map((table) => <div key={table.id} className="min-w-[112px] border-2 border-black bg-white px-3 py-2 font-mono text-[10px]"><p className="font-bold">{table.id}</p><p className="mt-1 text-black/60">{table.size}</p></div>)}</div></div>)}</div></div></Panel> }
