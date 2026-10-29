import { api } from "./api.js";

// Live copy of every job. The server pushes changes over server-sent
// events; on every (re)connect the full list is loaded again.
const jobs = new Map();
const listeners = new Set();

export const store = {
  /** list returns all jobs, newest first. */
  list: () => [...jobs.values()].sort((a, b) => Date.parse(b.createdAt) - Date.parse(a.createdAt)),
  /** subscribe calls fn after every change and returns an unsubscribe function. */
  subscribe(fn) {
    listeners.add(fn);
    return () => listeners.delete(fn);
  },
};

function emit() {
  for (const fn of listeners) fn();
}

function put(job) {
  if (job.status === "removed") jobs.delete(job.id);
  else jobs.set(job.id, job);
}

export function connect() {
  const source = new EventSource("/api/events");
  source.addEventListener("open", async () => {
    const list = await api.jobs().catch(() => null);
    if (!list) return;
    jobs.clear();
    list.forEach(put);
    emit();
  });
  source.addEventListener("job", (event) => {
    put(JSON.parse(event.data));
    emit();
  });
}
