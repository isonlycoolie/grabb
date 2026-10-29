import { api } from "../../api.js";
import { toast } from "../../lib/controls.js";
import { h, icon, siteMark } from "../../lib/dom.js";
import { store } from "../../store.js";
import { createJobList } from "../jobs/job-list.js";
import { createPreview, previewSkeleton } from "./preview.js";

const sites = ["youtube", "tiktok", "instagram", "twitter", "facebook", "vimeo", "soundcloud", "reddit"];

/** renderDownload is the home page: paste a link, pick a format, watch it run. */
export function renderDownload({ page }) {
  const input = h("input", {
    type: "url", name: "url", required: true, autocomplete: "off", spellcheck: "false",
    placeholder: "Paste a link from YouTube, TikTok, Instagram...", "aria-label": "Link",
  });
  const fetchButton = h("button", { class: "btn btn-primary btn-lg", type: "submit" }, "Fetch");
  const error = h("p", { class: "form-error", role: "alert", hidden: true });
  const preview = h("div");
  let request = 0;

  async function read() {
    const url = input.value.trim();
    if (!url) return input.focus();
    const id = ++request;
    error.hidden = true;
    fetchButton.disabled = true;
    preview.replaceChildren(previewSkeleton());
    try {
      const info = await api.probe(url);
      if (id === request) preview.replaceChildren(createPreview(url, info, queued));
    } catch (err) {
      if (id !== request) return;
      preview.replaceChildren();
      error.replaceChildren(icon("alert"), err.message);
      error.hidden = false;
    } finally {
      if (id === request) fetchButton.disabled = false;
    }
  }

  async function paste() {
    try {
      input.value = (await navigator.clipboard.readText()).trim();
      read();
    } catch {
      input.focus();
      toast("The clipboard is blocked here. Press Ctrl+V in the box instead.", "error");
    }
  }

  function queued(info) {
    request++;
    input.value = "";
    preview.replaceChildren();
    toast(`Downloading ${info.title || "your file"}`, "success");
  }

  input.addEventListener("paste", () => setTimeout(read));
  const form = h("form", { class: "link-form", onsubmit: (e) => (e.preventDefault(), read()) },
    h("label", { class: "input" }, icon("link"), input,
      h("button", { class: "btn btn-ghost btn-sm", type: "button", onclick: paste }, icon("clipboard", 14), "Paste")),
    fetchButton);

  const count = h("span", { class: "tag" });
  const activity = createJobList(h("div", { class: "empty" }, icon("download", 20),
    h("strong", {}, "Nothing downloading"), h("span", {}, "Paste a link above to start.")));
  const sync = () => {
    const n = activity.update(store.list().filter((job) => job.status !== "done"));
    count.textContent = n;
    count.hidden = n === 0;
  };

  page.append(h("div", { class: "page-inner" },
    h("section", { class: "hero" },
      h("img", { class: "hero-logo", src: "assets/logo.svg", alt: "grabb", width: 48, height: 48 }),
      h("h1", {}, "Paste a link. Get the file."),
      h("p", {}, "Video or audio from social networks in the best quality they serve, up to 8K and 60fps, saved straight to your computer."),
      form, error,
      h("div", { class: "sites" }, "Works with", sites.map((site) => h("span", { class: "tag" }, siteMark(site))), "and 1,800 more")),
    preview,
    h("section", {}, h("div", { class: "section-head" }, h("h2", {}, "Activity"), count), activity.el)));

  sync();
  input.focus();
  return store.subscribe(sync);
}
