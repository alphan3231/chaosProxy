'use client';

import { formatDistanceToNow } from 'date-fns';
import { Ghost, CheckCircle, AlertTriangle, Clock } from 'lucide-react';

interface Log {
    method: string;
    path: string;
    status: number;
    duration: string;
    timestamp: string;
}

export default function RecentRequests({ logs }: { logs: Log[] }) {
    if (!logs || logs.length === 0) {
        return (
            <div className="text-center p-8 text-neutral-500 italic bg-neutral-900/30 rounded-lg border border-neutral-800">
                Waiting for traffic...
            </div>
        );
    }

    return (
        <div className="overflow-hidden rounded-lg border border-neutral-800 bg-neutral-900/50">
            <table className="w-full text-sm text-left">
                <thead className="bg-neutral-900 text-neutral-400 font-medium">
                    <tr>
                        <th className="px-4 py-3">Method</th>
                        <th className="px-4 py-3">Path</th>
                        <th className="px-4 py-3">Status</th>
                        <th className="px-4 py-3">Duration</th>
                        <th className="px-4 py-3">Time</th>
                    </tr>
                </thead>
                <tbody className="divide-y divide-neutral-800 text-neutral-300">
                    {logs.map((log, i) => (
                        <tr key={i} className="hover:bg-neutral-800/50 transition-colors">
                            <td className="px-4 py-3 font-mono">
                                <span className={`px-2 py-1 rounded text-xs font-bold ${getMethodColor(log.method)}`}>
                                    {log.method}
                                </span>
                            </td>
                            <td className="px-4 py-3 font-mono truncate max-w-[200px]" title={log.path}>
                                {log.path}
                            </td>
                            <td className="px-4 py-3">
                                <StatusBadge status={log.status} />
                            </td>
                            <td className="px-4 py-3 font-mono text-neutral-500">{log.duration}</td>
                            <td className="px-4 py-3 text-neutral-500 whitespace-nowrap">
                                <div className="flex items-center gap-1">
                                    <Clock size={12} />
                                    {formatDistanceToNow(new Date(log.timestamp), { addSuffix: true })}
                                </div>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

function StatusBadge({ status }: { status: number }) {
    if (status >= 200 && status < 300) {
        return (
            <span className="flex items-center gap-1 text-green-400">
                <CheckCircle size={14} /> {status}
            </span>
        );
    }
    if (status >= 500) {
        return (
            <span className="flex items-center gap-1 text-red-500 font-bold">
                <AlertTriangle size={14} /> {status}
            </span>
        );
    }
    return <span className="text-neutral-400">{status}</span>;
}

function getMethodColor(method: string) {
    switch (method) {
        case 'GET': return 'bg-blue-900/30 text-blue-400';
        case 'POST': return 'bg-green-900/30 text-green-400';
        case 'PUT': return 'bg-yellow-900/30 text-yellow-400';
        case 'DELETE': return 'bg-red-900/30 text-red-400';
        default: return 'bg-neutral-800 text-neutral-400';
    }
}
