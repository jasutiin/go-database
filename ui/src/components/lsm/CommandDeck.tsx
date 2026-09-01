import { Database, Pencil, Play, Plus, Search, Trash2, Workflow } from 'lucide-react'
import { useState } from 'react'

import type { StorageOperation } from './types'

type CommandDeckProps = {
  onOperation: (operation: StorageOperation) => void
  onCompaction: () => void
}

const actions = [
  { kind: 'put' as const, icon: Plus },
  { kind: 'get' as const, icon: Search },
  { kind: 'update' as const, icon: Pencil },
  { kind: 'delete' as const, icon: Trash2 },
]

export function CommandDeck({ onOperation, onCompaction }: CommandDeckProps) {
  const [kind, setKind] = useState<StorageOperation['kind']>('put')
  const [key, setKey] = useState('user:73')
  const [value, setValue] = useState('active')
  const [error, setError] = useState('')
  const needsValue = kind === 'put' || kind === 'update'

  function enqueue() {
    if (!key.trim() || (needsValue && !value.trim())) {
      setError(needsValue ? 'key and value are required' : 'a key is required')
      return
    }
    setError('')
    onOperation({ kind, key: key.trim(), value: needsValue ? value : undefined })
  }

  return (
    <section className="border-2 border-black bg-black p-4 text-[#f8f8f2] shadow-[7px_7px_0_#111]">
      <div className="mb-4 border-b border-white/35 pb-3"><p className="font-mono text-[10px] font-bold uppercase tracking-[0.18em]">db-engine // command console</p><h2 className="mt-1 font-mono text-xl font-black uppercase">choose an operation</h2><p className="mt-2 font-mono text-[11px] text-white/65">the visualizer will replay each storage step as a single card.</p></div><div className="flex flex-col gap-3 xl:items-end">
        <div className="hidden min-w-[185px] xl:block"><p className="flex items-center gap-2 font-mono text-[10px] font-bold uppercase tracking-[0.18em]"><Database size={14} /> command deck</p></div>
        <div className="grid grid-cols-4 gap-1">
          {actions.map(({ kind: action, icon: Icon }) => <button type="button" key={action} onClick={() => setKind(action)} className={`flex items-center justify-center gap-1 border border-white px-2 py-2 font-mono text-[10px] font-bold uppercase ${kind === action ? 'bg-white text-black' : 'hover:bg-white hover:text-black'}`}><Icon size={13} />{action}</button>)}
        </div>
        <label className="flex-1 font-mono text-[10px] font-bold uppercase tracking-[0.12em]">key<input value={key} onChange={(event) => setKey(event.target.value)} className="mt-1 block w-full border border-white bg-white px-2 py-2 font-mono text-xs text-black outline-none" /></label>
        <label className="flex-1 font-mono text-[10px] font-bold uppercase tracking-[0.12em]">value<input value={value} disabled={!needsValue} onChange={(event) => setValue(event.target.value)} className="mt-1 block w-full border border-white bg-white px-2 py-2 font-mono text-xs text-black outline-none disabled:bg-white/30" /></label>
        <button type="button" onClick={enqueue} className="flex items-center justify-center gap-2 border-2 border-white bg-white px-4 py-2.5 font-mono text-xs font-bold uppercase text-black hover:bg-black hover:text-white"><Play size={14} fill="currentColor" /> enqueue</button>
        <button type="button" onClick={onCompaction} className="flex items-center justify-center gap-2 border border-white px-3 py-2.5 font-mono text-[10px] font-bold uppercase hover:bg-white hover:text-black"><Workflow size={14} /> compact</button>
      </div>
      {error && <p className="mt-2 font-mono text-[10px] font-bold uppercase text-white">// {error}</p>}
    </section>
  )
}
