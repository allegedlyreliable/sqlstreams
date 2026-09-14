import type { ReferenceSection } from './model';

export const sections: ReferenceSection[] = [
	{
		id: 'api',
		title: 'API',
		ids: ['reference/consumer-registration', 'reference/consume', 'reference/message-metadata'],
	},
	{ 
    id: 'cli', 
    title: 'CLI', 
    ids: ['reference/consumer-cli'] 
  },
	{
		id: 'configuration',
		title: 'Configuration',
		ids: ['reference/consumer-group-settings', 'reference/consumer-session-options'],
	},
	{
		id: 'metrics',
		title: 'Metrics',
		ids: ['reference/consumer-cursor-backlog'],
	},
	{ 
    id: 'logs', 
    title: 'Logs', 
    ids: ['reference/consumer-stopped-log'] 
  },
	{ 
    id: 'alerts', 
    title: 'Alerts', 
    ids: ['reference/worker-liveness-alert'] 
  },
];
