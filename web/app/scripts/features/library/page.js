import { segmented } from "../../lib/controls.js";
import { h, icon, pageHeader } from "../../lib/dom.js";
import * as fmt from "../../lib/format.js";
import { store } from "../../store.js";
import { createJobList } from "../jobs/job-list.js";

/** renderLibrary lists finished downloads with search and a type filter. */
export function renderLibrary({ page }) {
  let query = "";
  let kind = "all";
  const summary = h("span");
  const emptyTitle = h("strong");
  const emptyText = h("span");
  const list = createJobList(h("div", { class: "empty" }, icon("library", 20), emptyTitle, emptyText));

  const search = h("label", { class: "input" }, icon("search"), h("input", {
    type: "search", placeholder: "Search", "aria-label": "Search downloads",
    oninput: (e) => {
      query = e.target.value.trim().toLowerCase();
      sync();
    },
  }));
  const filter = segmented([["all", "All"], ["video", "Video", "video"], ["audio", "Audio", "music"]], kind, (value) => {
    kind = value;
    sync();
  }, "Type");

  function matches(job) {
    if (kind !== "all" && job.options.mode !== kind) return false;
    return !query || [job.title, job.uploader, job.site].join(" ").toLowerCase().includes(query);
  }

  function sync() {
    const done = store.list().filter((job) => job.status === "done");
    list.update(done.filter(matches));
    const total = done.reduce((sum, job) => sum + (job.size || 0), 0);
    summary.textContent = done.length
      ? `${done.length} ${done.length === 1 ? "file" : "files"} · ${fmt.bytes(total)}`
      : "Finished downloads show up here.";
    emptyTitle.textContent = done.length ? "No matches" : "No downloads yet";
    emptyText.textContent = done.length ? "Try another search or type." : "Paste a link on the Download page to get started.";
  }

  page.append(h("div", { class: "page-inner" },
    pageHeader("Library", summary, h("div", { class: "library-tools" }, search, filter)),
    list.el));
  sync();
  return store.subscribe(sync);
}
