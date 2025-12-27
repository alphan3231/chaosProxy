'use client';

import { useEffect, useState } from 'react';
import { Activity, Ghost, Server, ShieldAlert } from 'lucide-react';

interface Stats {
  totalRequests: number;
  ghostCount: number;
  status: string;
}

export default function Home() {
  const [stats, setStats] = useState<Stats | null>(null);

  const fetchStats = async () => {
    try {
      const res = await fetch('/api/stats');
      const data = await res.json();
      setStats(data);
    } catch (error) {
      console.error('Failed to fetch stats', error);
    }
  };

  useEffect(() => {
    fetchStats();
    const interval = setInterval(fetchStats, 2000);
    return () => clearInterval(interval);
  }, []);

  if (!stats) return <div className="flex h-screen items-center justify-center bg-black text-white">Loading Chaos...</div>;

  return (
    <main className="min-h-screen bg-neutral-950 text-neutral-100 p-8">
      <div className="max-w-6xl mx-auto space-y-8">

        {/* Header */}
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

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <StatCard
            title="Total Requests"
            value={stats.totalRequests.toLocaleString()}
            icon={<Activity className="w-6 h-6 text-blue-400" />}
            color="blue"
          />
          <StatCard
            title="Ghost Mode Activations"
            value={stats.ghostCount.toLocaleString()}
            icon={<ShieldAlert className="w-6 h-6 text-red-500 animate-pulse" />}
            color="red"
          />
          <StatCard
            title="Active Services"
            value="1"
            icon={<Server className="w-6 h-6 text-emerald-400" />}
            color="emerald"
          />
        </div>

        {/* Info Section */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="p-6 rounded-xl bg-neutral-900/50 border border-neutral-800 backdrop-blur-sm">
            <h3 className="text-xl font-semibold mb-4 flex items-center gap-2">
              <Ghost className="w-5 h-5" /> What is Ghost Mode?
            </h3>
            <p className="text-neutral-400 leading-relaxed">
              When your backend fails, Chaos Proxy automatically switches to Ghost Mode.
              It serves "learned" responses from Redis, ensuring your users never see an error page.
            </p>
          </div>

          <div className="p-6 rounded-xl bg-neutral-900/50 border border-neutral-800 backdrop-blur-sm">
            <h3 className="text-xl font-semibold mb-4 text-neutral-200">Recent Events</h3>
            <div className="space-y-3">
              <div className="text-sm text-neutral-500 italic">No recent anomalies detected...</div>
            </div>
          </div>
        </div>

      </div>
    </main>
  );
}

function StatCard({ title, value, icon, color }: { title: string, value: string, icon: any, color: string }) {
  return (
    <div className={`p-6 rounded-2xl bg-neutral-900/40 border border-neutral-800 hover:border-${color}-500/50 transition-all duration-300 group`}>
      <div className="flex items-start justify-between">
        <div>
          <p className="text-neutral-400 text-sm font-medium mb-1">{title}</p>
          <h2 className="text-4xl font-bold tracking-tight text-white group-hover:text-neutral-100 transition-colors">
            {value}
          </h2>
        </div>
        <div className={`p-3 rounded-lg bg-${color}-500/10 border border-${color}-500/20`}>
          {icon}
        </div>
      </div>
    </div>
  );
}
