// Client for the grabb server. Every call rejects with a readable message.
async function request(method, path, body) {
  let res;
  try {
    res = await fetch(`/api${path}`, {
      method,
      headers: body ? { "Content-Type": "application/json" } : {},
      body: body && JSON.stringify(body),
    });
  } catch {
    throw new Error("Can't reach grabb. Make sure it is still running, then try again.");
  }
  const data = res.status === 204 ? null : await res.json().catch(() => null);
  if (!res.ok) throw new Error(data?.error || `Request failed (${res.status})`);
  return data;
}

let settings;

export const api = {
  settings: () => (settings ??= request("GET", "/settings")),
  probe: (url) => request("POST", "/probe", { url }),
  jobs: () => request("GET", "/jobs"),
  start: (url, options, known) => request("POST", "/jobs", { url, ...options, known }),
  cancel: (id) => request("POST", `/jobs/${id}/cancel`),
  retry: (id) => request("POST", `/jobs/${id}/retry`),
  reveal: (id) => request("POST", `/jobs/${id}/reveal`),
  remove: (id) => request("DELETE", `/jobs/${id}`),
  fileURL: (id, inline = false) => `/api/jobs/${id}/file${inline ? "?inline=1" : ""}`,
};
