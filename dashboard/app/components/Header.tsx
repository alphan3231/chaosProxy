'use client';

import { Ghost } from 'lucide-react';

export default function Header() {
    return (
        <header className="flex items-center justify-between border-b border-neutral-800 pb-6">
            <div className="flex items-center space-x-3">
                <div className="p-3 bg-purple-900/20 rounded-lg border border-purple-500/30">
                    <Ghost className="w-8 h-8 text-purple-400" />
                </div>
                <div>
                    <h1 className="text-3xl font-bold bg-gradient-to-r from-purple-400 to-pink-400 bg-clip-text text-transparent">
                        Chaos Proxy
                    </h1>
                    <p className="text-neutral-500">Immortality Layer Dashboard</p>
                </div>
            </div>
            <div className="flex items-center space-x-2 text-sm">
                <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
                <span className="text-green-500">System Operational</span>
            </div>
        </header>
    );
}
