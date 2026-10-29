const units = ["B", "KB", "MB", "GB", "TB"];

/** bytes formats a size: 1536 -> "1.5 KB". */
export function bytes(n) {
  if (!(n > 0)) return "";
  let i = 0;
  while (n >= 1024 && i < units.length - 1) {
    n /= 1024;
    i++;
  }
  return `${n.toFixed(n >= 100 || i === 0 ? 0 : 1)} ${units[i]}`;
}

/** speed formats bytes per second. */
export const speed = (n) => (n > 0 ? `${bytes(n)}/s` : "");

/** clock formats seconds as 1:05 or 1:02:05. */
export function clock(seconds) {
  if (!(seconds >= 0)) return "";
  const s = Math.round(seconds);
  const hours = Math.floor(s / 3600);
  const minutes = Math.floor((s % 3600) / 60);
  const rest = String(s % 60).padStart(2, "0");
  return hours ? `${hours}:${String(minutes).padStart(2, "0")}:${rest}` : `${minutes}:${rest}`;
}

const relative = new Intl.RelativeTimeFormat(undefined, { numeric: "auto" });
const steps = [
  ["year", 31536000],
  ["month", 2592000],
  ["day", 86400],
  ["hour", 3600],
  ["minute", 60],
];

/** ago formats a timestamp relative to now: "5 minutes ago". */
export function ago(iso) {
  if (!iso) return "";
  const seconds = (new Date(iso).getTime() - Date.now()) / 1000;
  for (const [unit, size] of steps) {
    if (Math.abs(seconds) >= size) return relative.format(Math.round(seconds / size), unit);
  }
  return "just now";
}

/** resolution names a video height the way people say it. */
export function resolution(height) {
  if (height >= 4320) return "8K";
  if (height >= 2160) return "4K";
  return `${height}p`;
}

/** quality describes what a job was asked for. */
export function quality(options) {
  if (options.mode === "audio") return options.audioFormat.toUpperCase();
  const label = options.height ? resolution(options.height) : "Best quality";
  return options.compatible ? `${label} H.264` : label;
}

/** percent is how far the current stream has come, 0 to 100. */
export const percent = (p) => (p?.total > 0 ? Math.min(100, (p.downloaded / p.total) * 100) : 0);
