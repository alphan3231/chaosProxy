'use client';

import { useState, useEffect } from 'react';
import { ShieldBan, RefreshCw } from 'lucide-react';

export default function BlockedIPs() {
    const [ips, setIps] = useState<string[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchIPs = () => {
        setLoading(true);
        fetch('/api/blocked-ips')
            .then((res) => res.json())
            .then((data) => {
                setIps(data.blocked_ips || []);
                setLoading(false);
            })
            .catch((err) => {
                console.error(err);
                setLoading(false);
            });
    };

    useEffect(() => {
        fetchIPs();
    }, []);

    return (
        <div className="bg-[#1a1a1a] p-6 rounded-xl border border-gray-800 shadow-lg h-full">
            {/* Header */}
            <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-2 text-red-500">
                    <ShieldBan size={24} />
                    <h2 className="text-xl font-bold">Blocked IPs</h2>
                </div>
                <button onClick={fetchIPs} className="p-2 hover:bg-white/5 rounded-full text-gray-400 hover:text-white transition-all">
                    <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                </button>
            </div>

            {/* List */}
            <div className="space-y-2 max-h-[300px] overflow-y-auto pr-2 custom-scrollbar">
                {ips.length === 0 ? (
                    <div className="text-gray-500 text-sm italic py-4 text-center">No IPs are currently blocked.</div>
                ) : (
                    ips.map((ip, i) => (
                        <div key={i} className="flex items-center justify-between bg-black/40 p-3 rounded border border-gray-800 text-sm hover:border-red-900/50 transition-colors">
                            <span className="font-mono text-gray-300">{ip}</span>
                            <span className="text-xs text-red-500/80 uppercase tracking-wider font-bold">Blocked</span>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
}
