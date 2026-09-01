import type { LsmDashboardSnapshot } from './types'

export const mockLsmSnapshot: LsmDashboardSnapshot = {
  wal: {
    segment: 'wal-000014.log',
    entries: 48,
    usedBytes: '3.8 kb',
    capacity: '16 mb',
    status: 'live',
    records: ['put user:42', 'put session:81', 'delete cache:18', 'put user:73'],
  },
  activeMemtable: {
    entries: 48,
    usedBytes: '3.8 kb',
    limit: '64 kb',
    writeRate: '12 ops/s',
    status: 'live',
  },
  skipList: {
    levels: [
      ['user:42'],
      ['cache:18', 'user:42'],
      ['session:81', 'cache:18', 'user:42'],
      ['user:12', 'session:81', 'cache:18', 'user:42', 'user:73'],
    ],
    lookupSteps: 4,
  },
  immutableMemtables: [
    { id: 'mem-000013', entries: 256, size: '21.4 kb', status: 'flushing' },
    { id: 'mem-000012', entries: 256, size: '20.8 kb', status: 'queued' },
  ],
  compaction: {
    status: 'idle',
    inputTables: 0,
    outputTables: 0,
    debt: '0 bytes',
    lastRun: '12:04:19',
  },
  sstableLevels: [
    {
      level: 'l0',
      tables: [
        { id: '000009.sst', size: '21 kb' },
        { id: '000010.sst', size: '24 kb' },
      ],
    },
    {
      level: 'l1',
      tables: [
        { id: '000004.sst', size: '128 kb' },
        { id: '000006.sst', size: '117 kb' },
        { id: '000008.sst', size: '131 kb' },
      ],
    },
    { level: 'l2', tables: [{ id: '000001.sst', size: '512 kb' }] },
  ],
}
