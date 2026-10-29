import { api } from "../../api.js";
import { toast } from "../../lib/controls.js";
import { h, icon, siteMark, thumbnail } from "../../lib/dom.js";
import * as fmt from "../../lib/format.js";
import { openPlayer } from "../player/player.js";

/** createJobRow renders one job and updates it in place as it changes. */
export function createJobRow(job) {
  const thumb = h("button", { class: "job-thumb", type: "button", onclick: () => openPlayer(job) });
  const title = h("div", { class: "job-title" });
  const meta = h("div", { class: "job-meta" });
  const fill = h("div", { class: "progress-fill" });
  const bar = h("div", { class: "progress", role: "progressbar", "aria-label": "Progress" }, fill);
  const actions = h("div", { class: "job-actions" });
  const el = h("div", { class: "job-row" }, thumb, h("div", { class: "job-body" }, title, meta, bar), actions);
  let shownStatus;
  let shownThumb;

  function update(next) {
    job = next;
    const name = job.title || job.url;
    el.dataset.status = job.status;
    title.textContent = name;
    title.title = name;
    meta.replaceChildren(...details(job).filter(Boolean).map((part) => (part instanceof Node ? part : h("span", {}, part))));

    const pct = fmt.percent(job.progress);
    bar.hidden = job.status !== "running" && job.status !== "queued";
    bar.setAttribute("aria-valuenow", Math.round(pct));
    fill.style.width = `${pct}%`;

    if (job.thumbnail !== shownThumb) {
      shownThumb = job.thumbnail;
      thumb.replaceChildren(thumbnail(job.thumbnail, job.options.mode === "audio" ? "music" : "video"),
        h("span", { class: "play" }, icon("play", 14)));
    }
    if (job.status !== shownStatus) {
      shownStatus = job.status;
      actions.replaceChildren(...buttons(job));
      thumb.disabled = job.status !== "done";
      thumb.setAttribute("aria-label", job.status === "done" ? `Play ${name}` : name);
    }
  }

  update(job);
  return { el, update };
}

const tag = (text, color) => h("span", { class: `tag tag-${color}` }, text);

function details(job) {
  const p = job.progress ?? {};
  switch (job.status) {
    case "queued":
      return [tag("Queued", "gray"), job.site && siteMark(job.site)];
    case "running":
      return [
        tag(job.stage || "Starting", "blue"),
        p.total > 0 && `${Math.floor(fmt.percent(p))}%`,
        p.total > 0 ? `${fmt.bytes(p.downloaded)} of ${fmt.bytes(p.total)}` : fmt.bytes(p.downloaded),
        fmt.speed(p.speed),
        p.eta > 0 && `${fmt.clock(p.eta)} left`,
      ];
    case "done":
      return [job.site && siteMark(job.site), job.uploader, fmt.quality(job.options), fmt.bytes(job.size), fmt.ago(job.finishedAt)];
    case "failed":
      return [tag("Failed", "red"), h("span", { class: "job-error" }, job.error || "Download failed")];
    case "canceled":
      return [tag("Canceled", "gray")];
    default:
      return [];
  }
}

function button(label, iconName, action, variant = "") {
  const onclick = () => action().catch((err) => toast(err.message, "error"));
  return h("button", { class: `btn btn-ghost btn-sm ${variant}`, type: "button", title: label, onclick },
    icon(iconName, 14), h("span", {}, label));
}

function buttons(job) {
  const remove = button("Remove", "trash", () => api.remove(job.id), "btn-danger");
  switch (job.status) {
    case "queued":
    case "running":
      return [button("Cancel", "x", () => api.cancel(job.id))];
    case "done":
      return [
        button("Play", "play", async () => openPlayer(job)),
        button("Show in folder", "folder", () => api.reveal(job.id)),
        remove,
      ];
    default:
      return [button("Retry", "refresh", () => api.retry(job.id)), remove];
  }
}
