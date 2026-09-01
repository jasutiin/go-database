import type { LsmDashboardSnapshot, StorageOperation, TraceEvent } from './types'

export function createOperationTrace(operation: StorageOperation, sequence: number): TraceEvent[] {
  const event = (kind: string, target: TraceEvent['target'], title: string, detail: string): TraceEvent => ({
    id: `${sequence}-${kind}`,
    source: 'user',
    kind,
    target,
    title,
    detail,
    key: operation.key,
    sequence,
  })

  if (operation.kind === 'get') {
    return [
      event('lookup_started', 'skip-list', 'searching skip list', `locating ${operation.key} in the active memtable`),
      event('lookup_complete', 'memtable', 'lookup complete', `${operation.key} was checked in the active memtable`),
    ]
  }

  const verb = operation.kind === 'delete' ? 'delete' : 'put'
  const trace = [
    event('wal_appended', 'wal', 'wal append', `${verb} record for ${operation.key} added to the log`),
    event('wal_synced', 'wal', 'wal sync', 'the write-ahead log is durable before the memtable changes'),
    event(
      operation.kind === 'delete' ? 'skip_list_deleted' : 'skip_list_updated',
      'skip-list',
      operation.kind === 'delete' ? 'tombstone inserted' : 'skip list update',
      operation.kind === 'delete'
        ? `${operation.key} is marked as deleted in the ordered in-memory index`
        : `${operation.key} is inserted into the ordered in-memory index`,
    ),
    event('memtable_updated', 'memtable', 'memtable update', `${operation.key} is now visible to reads from the active memtable`),
  ]

  if (sequence % 3 === 0) {
    trace.push(
      event('flush_started', 'immutable', 'flush threshold reached', 'the active memtable becomes immutable and waits for disk'),
      event('sstable_created', 'sstables', 'sstable written', 'sorted immutable entries are installed in level zero'),
      event('flush_finished', 'memtable', 'new active memtable', 'a fresh skip list is ready for new writes'),
    )
  }

  return trace
}

export function createCompactionTrace(sequence: number): TraceEvent[] {
  const event = (kind: string, target: TraceEvent['target'], title: string, detail: string): TraceEvent => ({
    id: `${sequence}-${kind}`,
    source: 'background',
    kind,
    target,
    title,
    detail,
    sequence,
  })

  return [
    event('compaction_started', 'compaction', 'background compaction', 'level-zero tables have accumulated and are being merged'),
    event('compaction_merged', 'sstables', 'merge output written', 'the compacted table is placed into level one'),
    event('compaction_finished', 'compaction', 'compaction complete', 'obsolete versions can now be reclaimed'),
  ]
}

export function applyTraceEvent(snapshot: LsmDashboardSnapshot, event: TraceEvent): LsmDashboardSnapshot {
  switch (event.kind) {
    case 'wal_appended': {
      const entries = snapshot.wal.entries + 1
      return { ...snapshot, wal: { ...snapshot.wal, entries, usedBytes: `${(entries * 0.08).toFixed(1)} kb`, records: [`${event.detail}`, ...snapshot.wal.records].slice(0, 4) } }
    }
    case 'lookup_started':
    case 'lookup_complete':
      return { ...snapshot, skipList: { ...snapshot.skipList, lookupSteps: 2 + (event.sequence % 5) }, activeMemtable: { ...snapshot.activeMemtable, writeRate: `lookup / ${2 + (event.sequence % 5)} steps` } }
    case 'skip_list_updated':
    case 'skip_list_deleted': {
      const bottom = snapshot.skipList.levels.at(-1) ?? []
      const isDelete = event.kind === 'skip_list_deleted'
      const nextKey = isDelete ? `† ${event.key}` : event.key ?? 'key'
      const keys = [...bottom.filter((key) => key !== event.key && key !== `† ${event.key}`), nextKey].slice(-5)
      return { ...snapshot, skipList: { ...snapshot.skipList, levels: buildLevels(keys) } }
    }
    case 'memtable_updated': {
      const entries = snapshot.activeMemtable.entries + 1
      return { ...snapshot, activeMemtable: { ...snapshot.activeMemtable, entries, usedBytes: `${(entries * 0.08).toFixed(1)} kb`, writeRate: `${8 + (event.sequence % 11)} ops/s` } }
    }
    case 'flush_started':
      return { ...snapshot, immutableMemtables: [{ id: `mem-${String(event.sequence).padStart(6, '0')}`, entries: snapshot.activeMemtable.entries, size: snapshot.activeMemtable.usedBytes, status: 'flushing' }, ...snapshot.immutableMemtables].slice(0, 2) }
    case 'sstable_created':
      return { ...snapshot, sstableLevels: snapshot.sstableLevels.map((level, index) => index === 0 ? { ...level, tables: [{ id: `${String(event.sequence).padStart(6, '0')}.sst`, size: '24 kb' }, ...level.tables].slice(0, 4) } : level) }
    case 'flush_finished':
      return { ...snapshot, activeMemtable: { ...snapshot.activeMemtable, entries: 0, usedBytes: '0.0 kb' }, skipList: { ...snapshot.skipList, levels: [[], [], [], []] }, immutableMemtables: snapshot.immutableMemtables.map((table) => ({ ...table, status: 'queued' })) }
    case 'compaction_started':
      return { ...snapshot, compaction: { ...snapshot.compaction, status: 'flushing', inputTables: 2, outputTables: 1, debt: '45 kb' } }
    case 'compaction_merged':
      return { ...snapshot, sstableLevels: snapshot.sstableLevels.map((level, index) => index === 1 ? { ...level, tables: [{ id: `${String(event.sequence).padStart(6, '0')}-merged.sst`, size: '48 kb' }, ...level.tables] } : level) }
    case 'compaction_finished':
      return { ...snapshot, compaction: { ...snapshot.compaction, status: 'idle', inputTables: 0, outputTables: 0, debt: '0 bytes', lastRun: `seq ${event.sequence}` } }
    default:
      return snapshot
  }
}

function buildLevels(keys: string[]) {
  return [keys.filter((_, index) => index % 4 === 0), keys.filter((_, index) => index % 3 === 0), keys.filter((_, index) => index % 2 === 0), keys]
}
