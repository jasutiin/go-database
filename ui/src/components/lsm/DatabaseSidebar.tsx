import { Circle, Pencil, Plus, Search, Terminal, Trash2 } from 'lucide-react'
import { useState } from 'react'

import type { StorageOperation } from './types'

type DatabaseSidebarProps = {
  lastOperation: StorageOperation | null
  onOperation: (operation: StorageOperation) => void
}

const operationOptions: { kind: StorageOperation['kind']; label: string; icon: typeof Plus }[] = [
  { kind: 'put', label: 'put', icon: Plus },
  { kind: 'update', label: 'update', icon: Pencil },
  { kind: 'get', label: 'get', icon: Search },
  { kind: 'delete', label: 'delete', icon: Trash2 },
]

export function DatabaseSidebar({ lastOperation, onOperation }: DatabaseSidebarProps) {
  const [kind, setKind] = useState<StorageOperation['kind']>('put')
  const [key, setKey] = useState('user:73')
  const [value, setValue] = useState('active')
  const [error, setError] = useState<string | null>(null)
  const needsValue = kind === 'put' || kind === 'update'

  function submitOperation() {
    const trimmedKey = key.trim()
    if (!trimmedKey) {
      setError('a key is required')
      return
    }
    if (needsValue && !value.trim()) {
      setError('a value is required')
      return
    }

    setError(null)
    onOperation({ kind, key: trimmedKey, value: needsValue ? value : undefined })
  }

  return (
    <aside className="h-fit border-2 border-black bg-[#f8f8f2] shadow-[5px_5px_0_#111] lg:sticky lg:top-5">
      <header className="border-b-2 border-black bg-black px-3 py-3 text-[#f8f8f2]">
        <p className="flex items-center gap-2 font-mono text-[10px] font-bold uppercase tracking-[0.18em]"><Terminal size={14} /> database console</p>
        <p className="mt-1 font-mono text-[10px] uppercase text-white/60">lsm // local interface</p>
      </header>

      <div className="p-3">
        <p className="mb-2 font-mono text-[10px] font-bold uppercase tracking-[0.14em]">operation</p>
        <div className="grid grid-cols-2 gap-2">
          {operationOptions.map(({ kind: optionKind, label, icon: Icon }) => (
            <button
              type="button"
              key={optionKind}
              onClick={() => setKind(optionKind)}
              className={`flex items-center justify-between border-2 border-black px-2 py-2 font-mono text-[11px] font-bold uppercase transition ${kind === optionKind ? 'bg-black text-white' : 'bg-white hover:bg-black hover:text-white'}`}
            >
              {label}
              <Icon size={14} strokeWidth={2.5} />
            </button>
          ))}
        </div>
      </div>

      <div className="border-y-2 border-black bg-white p-3">
        <label className="block font-mono text-[10px] font-bold uppercase tracking-[0.14em]" htmlFor="storage-key">key</label>
        <input
          id="storage-key"
          value={key}
          onChange={(event) => setKey(event.target.value)}
          placeholder="user:42"
          className="mt-1.5 w-full border-2 border-black bg-[#f8f8f2] px-2 py-2 font-mono text-xs outline-none placeholder:text-black/35 focus:bg-black focus:text-white"
        />

        <label className="mt-4 block font-mono text-[10px] font-bold uppercase tracking-[0.14em]" htmlFor="storage-value">value</label>
        <textarea
          id="storage-value"
          value={value}
          onChange={(event) => setValue(event.target.value)}
          placeholder={needsValue ? 'value' : 'not used for this operation'}
          disabled={!needsValue}
          rows={4}
          className="mt-1.5 w-full resize-none border-2 border-black bg-[#f8f8f2] px-2 py-2 font-mono text-xs outline-none placeholder:text-black/35 disabled:cursor-not-allowed disabled:bg-black/10 focus:bg-black focus:text-white"
        />

        {error && <p className="mt-2 border border-black bg-black px-2 py-1.5 font-mono text-[10px] uppercase text-white">{error}</p>}
      </div>

      <div className="p-3">
        <button type="button" onClick={submitOperation} className="flex w-full items-center justify-center gap-2 border-2 border-black bg-black px-3 py-3 font-mono text-xs font-bold uppercase tracking-[0.16em] text-white transition hover:bg-white hover:text-black">
          run {kind}
        </button>
        <p className="mt-3 font-mono text-[10px] leading-relaxed text-black/60">controls currently capture local intent only. connect this action to the storage api when the server mutation endpoint exists.</p>
      </div>

      <footer className="border-t-2 border-black bg-white px-3 py-2">
        <p className="flex items-center gap-2 font-mono text-[10px] font-bold uppercase tracking-[0.12em]"><Circle size={9} fill="currentColor" /> {lastOperation ? `queued ${lastOperation.kind}` : 'awaiting operation'}</p>
        {lastOperation && <p className="mt-1 truncate font-mono text-[10px] text-black/60">{lastOperation.key}{lastOperation.value ? ` = ${lastOperation.value}` : ''}</p>}
      </footer>
    </aside>
  )
}
