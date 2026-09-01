import { createFileRoute } from '@tanstack/react-router'

import { LsmDashboard } from '../components/lsm/LsmDashboard'

export const Route = createFileRoute('/')({ component: LsmDashboard })
