import React, { useEffect, useMemo, useState } from "react";
import { api, Habit, Log } from "../lib/api";

type ModalProps = {
	open: boolean;
	onClose: () => void;
	children: React.ReactNode;
};

function Modal({ open, onClose, children }: ModalProps) {
	if (!open) return null;
	return (
		<div className="modal-backdrop" onClick={onClose}>
			<div className="modal" onClick={(e) => e.stopPropagation()}>
				{children}
			</div>
		</div>
	);
}

export function HabitUI() {
	const [habits, setHabits] = useState<Habit[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const [editingHabit, setEditingHabit] = useState<Habit | null>(null);
	const [editingLog, setEditingLog] = useState<{ habitId: number; log?: Log } | null>(null);
	const [refreshSignal, setRefreshSignal] = useState(0);

	useEffect(() => {
		let cancelled = false;
		setLoading(true);
		api
			.getHabits()
			.then((h) => {
				if (!cancelled) setHabits(h);
			})
			.catch((e) => !cancelled && setError(e.message))
			.finally(() => !cancelled && setLoading(false));
		return () => {
			cancelled = true;
		};
	}, []);

	const streaks = useMemo(() => {
		return new Map<number, number>();
	}, [habits]);

	return (
		<div className="container">
			<div className="header">
				<h1>Habits</h1>
				<button className="primary" onClick={() => setEditingHabit({ id: 0, user_id: 0, name: "", description: "", created_at: new Date().toISOString() })}>Add New Habit</button>
			</div>
			{loading && <p>Loading…</p>}
			{error && <p className="error">{error}</p>}
			<div className="habit-list">
				{habits.map((h) => (
					<HabitItem
						key={h.id}
						habit={h}
						streakCount={streaks.get(h.id) || 0}
						onEdit={() => setEditingHabit(h)}
						onDelete={async () => {
							if (!confirm("Delete this habit?")) return;
							await api.deleteHabit(h.id);
							setHabits((prev) => prev.filter((x) => x.id !== h.id));
						}}
						onNewLog={() => setEditingLog({ habitId: h.id })}
						onEditLog={(log) => setEditingLog({ habitId: h.id, log })}
						refreshSignal={refreshSignal}
					/>
				))}
			</div>

			<Modal open={!!editingHabit} onClose={() => setEditingHabit(null)}>
				<HabitForm
					initial={editingHabit || undefined}
					onCancel={() => setEditingHabit(null)}
					onSave={async (values) => {
						if (!editingHabit || editingHabit.id === 0) {
							const created = await api.createHabit(values);
							setHabits((prev) => [created, ...prev]);
						} else {
							const updated = await api.updateHabit(editingHabit.id, values);
							setHabits((prev) => prev.map((x) => (x.id === updated.id ? updated : x)));
						}
						setEditingHabit(null);
					}}
				/>
			</Modal>

			<Modal open={!!editingLog} onClose={() => setEditingLog(null)}>
				<LogForm
					initial={editingLog?.log}
					onCancel={() => setEditingLog(null)}
					onSave={async (values) => {
						if (!editingLog) return;
						if (editingLog.log) {
							await api.updateLog(editingLog.log.id, values);
						} else {
							await api.createLog(editingLog.habitId, values);
						}
						setEditingLog(null);
						setRefreshSignal((x) => x + 1);
					}}
				/>
			</Modal>
		</div>
	);
}

function HabitItem({ habit, streakCount, onEdit, onDelete, onNewLog, onEditLog, refreshSignal }: { habit: Habit; streakCount: number; onEdit: () => void; onDelete: () => void; onNewLog: () => void; onEditLog: (log: Log) => void; refreshSignal: number }) {
	const [logs, setLogs] = useState<Log[]>([]);
	const [loadingLogs, setLoadingLogs] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		setLoadingLogs(true);
		api
			.getHabitLogs(habit.id)
			.then((l) => {
				if (!cancelled) setLogs(l);
			})
			.catch((e) => !cancelled && setError(e.message))
			.finally(() => !cancelled && setLoadingLogs(false));
		return () => {
			cancelled = true;
		};
	// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [habit.id, refreshSignal]);

	const { gridDays, streak } = useMemo(() => {
		const days = buildRecentDays(30);
		const byDay = new Map<string, Log>();
		for (const log of logs) {
			const d = toLocalDate(log.created_at);
			if (!byDay.has(d)) byDay.set(d, log);
		}
		const filled = days.map((d) => ({ date: d, log: byDay.get(d) }));
		const s = computeStreakWithOneSkip(filled.map((x) => !!x.log));
		return { gridDays: filled, streak: s };
	}, [logs]);
	return (
		<div className="habit">
			<div className="habit-main">
				<div>
					<div className="habit-name">{habit.name}</div>
					<div className="habit-desc">{habit.description}</div>
				</div>
				<div className="habit-actions">
					<span className="streak">🔥 {streak}</span>
					<button onClick={onNewLog}>New Log</button>
					<button onClick={onEdit}>Edit</button>
					<button className="danger" onClick={onDelete}>Delete</button>
				</div>
			</div>
			<div className="streak-grid">
				{gridDays.map((d, idx) => (
					<div
						key={d.date}
						className={"day" + (d.log ? " filled" : "")}
						title={d.date}
						onClick={() => d.log && onEditLog(d.log)}
					/>
				))}
			</div>
		</div>
	);
}

function HabitForm({ initial, onCancel, onSave }: { initial?: Habit | null; onCancel: () => void; onSave: (values: { name?: string; description?: string; duration_seconds?: number }) => Promise<void> }) {
	const [name, setName] = useState(initial?.name || "");
	const [description, setDescription] = useState(initial?.description || "");
	const [duration, setDuration] = useState<string>(initial?.duration_seconds ? String(initial.duration_seconds) : "");
	const [error, setError] = useState<string | null>(null);
	const [saving, setSaving] = useState(false);

	return (
		<div>
			<h2>{initial && initial.id !== 0 ? "Edit Habit" : "New Habit"}</h2>
			{error && <p className="error">{error}</p>}
			<div className="form">
				<label>
					Name
					<input value={name} onChange={(e) => setName(e.target.value)} placeholder="Name" />
				</label>
				<label>
					Description
					<textarea value={description} onChange={(e) => setDescription(e.target.value)} placeholder="Description" />
				</label>
				<label>
					Duration seconds (optional)
					<input value={duration} onChange={(e) => setDuration(e.target.value)} inputMode="numeric" pattern="[0-9]*" placeholder="e.g. 1800" />
				</label>
			</div>
			<div className="modal-actions">
				<button onClick={onCancel} className="ghost">Cancel</button>
				<button
					className="primary"
					disabled={saving || !name.trim()}
					onClick={async () => {
						setSaving(true);
						setError(null);
						try {
							const durationVal = duration.trim() ? Number(duration) : undefined;
							if (durationVal !== undefined && Number.isNaN(durationVal)) throw new Error("Duration must be a number");
							await onSave({ name: name.trim(), description: description.trim() || undefined, duration_seconds: durationVal });
						} catch (e: any) {
							setError(e.message || "Failed to save");
						} finally {
							setSaving(false);
						}
					}}
				>
					Save
				</button>
			</div>
		</div>
	);
}

function toLocalDate(iso: string): string {
	const dt = new Date(iso);
	const y = dt.getFullYear();
	const m = String(dt.getMonth() + 1).padStart(2, "0");
	const d = String(dt.getDate()).padStart(2, "0");
	return `${y}-${m}-${d}`;
}

function buildRecentDays(n: number): string[] {
	const out: string[] = [];
	const now = new Date();
	for (let i = 0; i < n; i++) {
		const d = new Date(now);
		d.setDate(now.getDate() - (n - 1 - i));
		const y = d.getFullYear();
		const m = String(d.getMonth() + 1).padStart(2, "0");
		const dd = String(d.getDate()).padStart(2, "0");
		out.push(`${y}-${m}-${dd}`);
	}
	return out;
}

function computeStreakWithOneSkip(present: boolean[]): number {
	// Walk from the end (today) backwards, count consecutive days of presence
	// allowing one single skip day.
	let streak = 0;
	let skipsUsed = 0;
	for (let i = present.length - 1; i >= 0; i--) {
		if (present[i]) {
			streak++;
			continue;
		}
		if (skipsUsed === 0) {
			skipsUsed = 1;
			continue;
		}
		break;
	}
	return streak;
}

function LogForm({ initial, onCancel, onSave }: { initial?: Log; onCancel: () => void; onSave: (values: { notes?: string; duration_seconds?: number }) => Promise<void> }) {
	const [notes, setNotes] = useState(initial?.notes || "");
	const [duration, setDuration] = useState<string>(initial?.duration_seconds ? String(initial.duration_seconds) : "");
	const [error, setError] = useState<string | null>(null);
	const [saving, setSaving] = useState(false);

	return (
		<div>
			<h2>{initial ? "Edit Log" : "New Log"}</h2>
			{error && <p className="error">{error}</p>}
			<div className="form">
				<label>
					Notes
					<textarea value={notes} onChange={(e) => setNotes(e.target.value)} placeholder="Notes" />
				</label>
				<label>
					Duration seconds (optional)
					<input value={duration} onChange={(e) => setDuration(e.target.value)} inputMode="numeric" pattern="[0-9]*" placeholder="e.g. 600" />
				</label>
			</div>
			<div className="modal-actions">
				<button onClick={onCancel} className="ghost">Cancel</button>
				<button
					className="primary"
					disabled={saving}
					onClick={async () => {
						setSaving(true);
						setError(null);
						try {
							const durationVal = duration.trim() ? Number(duration) : undefined;
							if (durationVal !== undefined && Number.isNaN(durationVal)) throw new Error("Duration must be a number");
							await onSave({ notes: notes.trim() || undefined, duration_seconds: durationVal });
						} catch (e: any) {
							setError(e.message || "Failed to save");
						} finally {
							setSaving(false);
						}
					}}
				>
					Save
				</button>
			</div>
		</div>
	);
}


