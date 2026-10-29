import { api } from "../../api.js";
import { toast } from "../../lib/controls.js";
import { h, icon, siteMark, thumbnail } from "../../lib/dom.js";
import * as fmt from "../../lib/format.js";

/** openPlayer plays a finished download in a dialog over the page. */
export function openPlayer(job) {
  const src = api.fileURL(job.id, true);
  const audio = job.options.mode === "audio";
  const media = audio
    ? h("audio", { src, controls: true, autoplay: true, preload: "auto" })
    : h("video", { class: "player-video", src, controls: true, autoplay: true, playsinline: true, preload: "auto" });
  const stage = audio ? audioStage(job, media) : media;

  media.addEventListener("error", () => stage.replaceWith(h("p", { class: "player-note" }, icon("alert"),
    "Your browser can't play this file's format. Open it from its folder to watch it in your media player.")), { once: true });

  const dialog = h("dialog", { class: "player", "aria-label": job.title || "Player", autofocus: true },
    h("header", { class: "player-head" }, h("h2", {}, job.title || "Untitled"),
      h("button", { class: "btn btn-ghost icon-button", type: "button", "aria-label": "Close", onclick: close }, icon("x"))),
    stage,
    h("footer", { class: "player-foot" },
      job.site && siteMark(job.site), h("span", {}, fmt.quality(job.options)), h("span", {}, fmt.bytes(job.size)),
      h("div", { class: "actions" },
        h("button", { class: "btn btn-secondary btn-sm", type: "button", onclick: reveal }, icon("folder", 14), "Show in folder"),
        h("a", { class: "btn btn-secondary btn-sm", href: api.fileURL(job.id), download: "" }, icon("download", 14), "Save a copy"))));

  function reveal() {
    api.reveal(job.id).catch((err) => toast(err.message, "error"));
  }

  // Animate out before the dialog leaves the top layer.
  let closing = false;
  function close() {
    if (closing) return;
    closing = true;
    media.pause();
    dialog.classList.add("closing");
    dialog.addEventListener("animationend", () => dialog.close(), { once: true });
  }

  dialog.addEventListener("cancel", (event) => {
    event.preventDefault();
    close();
  });
  dialog.addEventListener("click", (event) => {
    if (event.target === dialog) close();
  });
  dialog.addEventListener("close", () => {
    removeEventListener("hashchange", close);
    dialog.remove();
  });
  addEventListener("hashchange", close);
  document.body.append(dialog);
  dialog.showModal();
}

function audioStage(job, media) {
  const backdrop = h("div", { class: "backdrop" });
  if (job.thumbnail) backdrop.style.backgroundImage = `url("${encodeURI(job.thumbnail)}")`;
  return h("div", { class: "player-audio" }, backdrop,
    h("div", { class: "art" }, thumbnail(job.thumbnail, "music")),
    job.uploader && h("p", {}, job.uploader), media);
}
