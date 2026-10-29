import { api } from "../../api.js";
import { chips, segmented, toast } from "../../lib/controls.js";
import { h, icon, siteMark, thumbnail } from "../../lib/dom.js";
import * as fmt from "../../lib/format.js";

const audioFormats = [
  ["mp3", "MP3"],
  ["m4a", "M4A"],
  ["opus", "Opus"],
  ["flac", "FLAC"],
];

/** createPreview shows what a link offers and lets the person pick a format. */
export function createPreview(url, info, onQueued) {
  const resolutions = (info.resolutions ?? []).filter((r) => r >= 240);
  const best = resolutions[0];
  const choice = { mode: best ? "video" : "audio", height: 0, audioFormat: "mp3", compatible: false };
  const options = h("div", { class: "preview-options" });
  const start = h("button", { class: "btn btn-primary btn-lg", type: "button", onclick: queue });

  function render() {
    const video = choice.mode === "video";
    const qualities = [[0, best ? `Best · ${fmt.resolution(best)}` : "Best"], ...resolutions.slice(1).map((r) => [r, fmt.resolution(r)])];
    options.replaceChildren(...[
      field("Format", segmented([["video", "Video", "video", !best], ["audio", "Audio", "music"]], choice.mode,
        (mode) => update({ mode }, true), "Format")),
      video
        ? field("Quality", chips(qualities, choice.height, (height) => update({ height }), "Quality"))
        : field("Audio format", chips(audioFormats, choice.audioFormat, (audioFormat) => update({ audioFormat }), "Audio format")),
      video && h("label", { class: "check" },
        h("input", { type: "checkbox", checked: choice.compatible, onchange: (e) => update({ compatible: e.target.checked }) }),
        h("span", {}, "Plays everywhere", h("small", {}, "H.264 MP4 for older phones, TVs and editors. May be below the best quality."))),
    ].filter(Boolean));
    label();
  }

  function update(patch, rerender = false) {
    Object.assign(choice, patch);
    if (rerender) render();
    else label();
  }

  function label() {
    const what = choice.mode === "audio" ? choice.audioFormat.toUpperCase() : fmt.resolution(choice.height || best);
    start.replaceChildren(icon("download"), `Download ${what}`);
  }

  async function queue() {
    start.disabled = true;
    const known = { title: info.title, uploader: info.uploader, extractor_key: info.site, thumbnail: info.thumbnail, duration: info.duration };
    try {
      await api.start(url, choice, known);
      onQueued(info);
    } catch (err) {
      toast(err.message, "error");
      start.disabled = false;
    }
  }

  render();
  return h("article", { class: "card preview-card" },
    h("div", { class: "preview-media" }, thumbnail(info.thumbnail, "video"),
      info.duration > 0 && h("span", { class: "tag" }, fmt.clock(info.duration))),
    h("div", { class: "preview-body" },
      h("div", {}, h("span", { class: "tag" }, siteMark(info.site)), h("h2", { class: "preview-title" }, info.title),
        info.uploader && h("div", { class: "preview-meta" }, info.uploader)),
      options,
      h("div", { class: "preview-actions" }, start, h("span", {}, "Saved to your grabb folder"))));
}

/** previewSkeleton holds the preview's place while the link is read. */
export function previewSkeleton() {
  return h("div", { class: "card preview-card", "aria-busy": "true", "aria-label": "Reading link" },
    h("div", { class: "preview-media skeleton" }),
    h("div", { class: "preview-body" }, h("div", { class: "skeleton line-sm" }), h("div", { class: "skeleton line-lg" }),
      h("div", { class: "skeleton line-md" })));
}

function field(label, control) {
  return h("div", { class: "field" }, h("span", { class: "field-label" }, label), control);
}
