export type PanelStatus = 'idle' | 'live' | 'queued' | 'flushing'

export type StorageOperation = {
  kind: 'put' | 'update' | 'get' | 'delete'
  key: string
  value?: string
}

export type WalSnapshot = {
  segment: string
  entries: number
  usedBytes: string
  capacity: string
  status: PanelStatus
  records: string[]
}

export type MemtableSnapshot = {
  entries: number
  usedBytes: string
  limit: string
  writeRate: string
  status: PanelStatus
}

export type SkipListSnapshot = {
  levels: string[][]
  lookupSteps: number
}

export type ImmutableMemtableSnapshot = {
  id: string
  entries: number
  size: string
  status: PanelStatus
}

export type CompactionSnapshot = {
  status: PanelStatus
  inputTables: number
  outputTables: number
  debt: string
  lastRun: string
}

export type SstableLevelSnapshot = {
  level: string
  tables: { id: string; size: string }[]
}


export type LsmDashboardSnapshot = {
  wal: WalSnapshot
  activeMemtable: MemtableSnapshot
  skipList: SkipListSnapshot
  immutableMemtables: ImmutableMemtableSnapshot[]
  compaction: CompactionSnapshot
  sstableLevels: SstableLevelSnapshot[]
}
