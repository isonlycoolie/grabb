import { h } from "../../lib/dom.js";
import { createJobRow } from "./job-row.js";

/**
 * createJobList keeps one row per job and reuses rows across updates, so
 * progress redraws never lose hover or focus. empty shows when nothing
 * matches.
 */
export function createJobList(empty) {
  const rows = new Map();
  const list = h("div", { class: "job-list" });
  const el = h("div", {}, list, empty);

  /** update renders jobs in order and returns how many are shown. */
  function update(jobs) {
    const ids = new Set(jobs.map((job) => job.id));
    for (const [id, row] of rows) {
      if (!ids.has(id)) {
        row.el.remove();
        rows.delete(id);
      }
    }
    jobs.forEach((job, i) => {
      let row = rows.get(job.id);
      if (row) row.update(job);
      else rows.set(job.id, (row = createJobRow(job)));
      if (list.children[i] !== row.el) list.insertBefore(row.el, list.children[i] ?? null);
    });
    list.hidden = jobs.length === 0;
    empty.hidden = jobs.length > 0;
    return jobs.length;
  }

  return { el, update };
}
