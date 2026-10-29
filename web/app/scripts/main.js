import { api } from "./api.js";
import { toast } from "./lib/controls.js";
import { icon } from "./lib/dom.js";
import { connect, store } from "./store.js";
import { renderDownload } from "./features/download/page.js";
import { renderLibrary } from "./features/library/page.js";
import { renderSettings } from "./features/settings/page.js";

// Each route renders into the page and may return a cleanup function.
const routes = {
  download: { title: "Download", icon: "download", render: renderDownload },
  library: { title: "Library", icon: "library", render: renderLibrary },
  settings: { title: "Settings", icon: "settings", render: renderSettings },
};

const page = document.getElementById("page");
const crumbs = document.getElementById("crumbs");
let cleanup;

function navigate() {
  const route = routes[location.hash.slice(2)] ?? routes.download;
  cleanup?.();
  page.replaceChildren();
  page.scrollTop = 0;
  crumbs.replaceChildren(icon(route.icon), route.title);
  document.title = `${route.title} · grabb`;
  for (const link of document.querySelectorAll("[data-route]")) {
    const current = routes[link.dataset.route] === route;
    link.classList.toggle("active", current);
    if (current) link.setAttribute("aria-current", "page");
    else link.removeAttribute("aria-current");
  }
  document.body.classList.remove("nav-open");
  cleanup = route.render({ page });
}

function decorate() {
  for (const el of document.querySelectorAll("[data-icon]")) el.prepend(icon(el.dataset.icon));
  document.getElementById("menu").addEventListener("click", () => document.body.classList.toggle("nav-open"));
  document.getElementById("scrim").addEventListener("click", () => document.body.classList.remove("nav-open"));
  api.settings().then((s) => (document.getElementById("version").textContent = `v${s.version}`), () => {});
}

// Counts the library in the sidebar and announces finished or failed jobs.
function watchJobs() {
  const count = document.getElementById("library-count");
  const seen = new Map();
  store.subscribe(() => {
    const jobs = store.list();
    count.textContent = jobs.filter((job) => job.status === "done").length || "";
    for (const job of jobs) {
      const before = seen.get(job.id);
      if (before && before !== job.status && job.status === "done") toast(`Saved ${job.title || "your download"}`, "success");
      if (before && before !== job.status && job.status === "failed") toast(job.error || "A download failed", "error");
      seen.set(job.id, job.status);
    }
  });
}

decorate();
watchJobs();
connect();
addEventListener("hashchange", navigate);
navigate();
