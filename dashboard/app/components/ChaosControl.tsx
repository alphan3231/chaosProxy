'use client';

import { useState, useEffect } from 'react';
import { Skull, AlertTriangle, Activity, Save } from 'lucide-react';

interface ChaosSettings {
    latency_enabled: boolean;
    latency_min: number;
    latency_max: number;
    failure_enabled: boolean;
    failure_rate: number;
}

export default function ChaosControl() {
    const [settings, setSettings] = useState<ChaosSettings>({
        latency_enabled: false,
        latency_min: 500,
        latency_max: 1500,
        failure_enabled: false,
        failure_rate: 10,
    });
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        fetch('/api/chaos')
            .then((res) => res.json())
            .then((data) => {
                setSettings(data);
                setLoading(false);
            })
            .catch((err) => console.error(err));
    }, []);

    const handleSave = async () => {
        setSaving(true);
        try {
            await fetch('/api/chaos', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Basic ' + btoa('admin:chaos123'), // Basic Auth Demo
                },
                body: JSON.stringify(settings),
            });
            alert('Chaos Settings Saved! 😈');
        } catch (error) {
            alert('Failed to save settings');
        } finally {
            setSaving(false);
        }
    };

    if (loading) return <div className="text-white">Loading Chaos Controls...</div>;

    return (
        <div className="bg-[#1a1a1a] p-6 rounded-xl border border-red-900/50 shadow-[0_0_20px_rgba(255,0,0,0.1)]">
            <div className="flex items-center gap-2 mb-6 text-red-500">
                <Skull size={24} />
                <h2 className="text-xl font-bold">Chaos Control</h2>
            </div>

            <div className="space-y-6">
                {/* Latency Control */}
                <div className="space-y-3">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2 text-yellow-500">
                            <Activity size={20} />
                            <label className="font-semibold">Latency Injection</label>
                        </div>
                        <input
                            type="checkbox"
                            checked={settings.latency_enabled}
                            onChange={(e) => setSettings({ ...settings, latency_enabled: e.target.checked })}
                            className="w-5 h-5 accent-red-600 cursor-pointer"
                        />
                    </div>

                    <div className={`space-y-2 pl-7 transition-opacity ${settings.latency_enabled ? 'opacity-100' : 'opacity-40 pointer-events-none'}`}>
                        <div className="flex gap-4">
                            <div className="flex-1">
                                <span className="text-xs text-gray-400">Min (ms)</span>
                                <input
                                    type="number"
                                    value={settings.latency_min}
                                    onChange={(e) => setSettings({ ...settings, latency_min: Number(e.target.value) })}
                                    className="w-full bg-black/50 border border-gray-700 rounded px-2 py-1 text-white text-sm"
                                />
                            </div>
                            <div className="flex-1">
                                <span className="text-xs text-gray-400">Max (ms)</span>
                                <input
                                    type="number"
                                    value={settings.latency_max}
                                    onChange={(e) => setSettings({ ...settings, latency_max: Number(e.target.value) })}
                                    className="w-full bg-black/50 border border-gray-700 rounded px-2 py-1 text-white text-sm"
                                />
                            </div>
                        </div>
                    </div>
                </div>

                <div className="h-px bg-gray-800" />

                {/* Failure Control */}
                <div className="space-y-3">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2 text-red-500">
                            <AlertTriangle size={20} />
                            <label className="font-semibold">Random Failure (500)</label>
                        </div>
                        <input
                            type="checkbox"
                            checked={settings.failure_enabled}
                            onChange={(e) => setSettings({ ...settings, failure_enabled: e.target.checked })}
                            className="w-5 h-5 accent-red-600 cursor-pointer"
                        />
                    </div>

                    <div className={`space-y-2 pl-7 transition-opacity ${settings.failure_enabled ? 'opacity-100' : 'opacity-40 pointer-events-none'}`}>
                        <div className="flex items-center justify-between text-sm text-gray-400">
                            <span>Failure Rate</span>
                            <span>{settings.failure_rate}%</span>
                        </div>
                        <input
                            type="range"
                            min="0"
                            max="100"
                            value={settings.failure_rate}
                            onChange={(e) => setSettings({ ...settings, failure_rate: Number(e.target.value) })}
                            className="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer accent-red-600"
                        />
                    </div>
                </div>

                <button
                    onClick={handleSave}
                    disabled={saving}
                    className="w-full flex items-center justify-center gap-2 bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-4 rounded transition-colors disabled:opacity-50"
                >
                    <Save size={18} />
                    {saving ? 'Applying Chaos...' : 'Apply Chaos Settings'}
                </button>
            </div>
        </div>
    );
}
