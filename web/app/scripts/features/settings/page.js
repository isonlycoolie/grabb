import { api } from "../../api.js";
import { segmented } from "../../lib/controls.js";
import { h, icon, pageHeader } from "../../lib/dom.js";

const themes = [
  ["system", "System"],
  ["light", "Light"],
  ["dark", "Dark"],
];

/** renderSettings shows how grabb is set up and the theme switch. */
export function renderSettings({ page }) {
  const body = h("div", { class: "page-inner" }, pageHeader("Settings", "How grabb saves files and what it runs on."));
  page.append(body);

  api.settings().then(
    (s) => body.append(
      group("Downloads", [
        ["Save folder", h("code", {}, s.outDir)],
        ["Downloads at once", s.workers],
        ["Connections per file", s.connections],
      ]),
      group("Appearance", [
        ["Theme", segmented(themes, window.grabbTheme.get(), (theme) => window.grabbTheme.set(theme), "Theme")],
      ]),
      group("Engine", Object.entries(s.tools).map(([name, version]) => [name, version
        ? h("code", {}, version)
        : h("span", { class: "missing" }, name === "aria2c" ? "Not installed. Files download over one connection." : "Not installed")])),
      group("About", [
        ["Version", s.version],
        ["License", "MIT"],
        ["Source", link("https://github.com/isonlycoolie/grabb", "github.com/isonlycoolie/grabb")],
        ["Built by", link("https://github.com/isonlycoolie", "isonlycoolie")],
      ]),
    ),
    (err) => body.append(h("p", { class: "form-error" }, icon("alert"), err.message)),
  );
}

function group(title, rows) {
  return h("section", { class: "card settings-group" }, h("h2", {}, title),
    h("dl", {}, rows.map(([label, value]) => h("div", { class: "setting" }, h("dt", {}, label), h("dd", {}, value)))));
}

function link(href, text) {
  return h("a", { href, target: "_blank", rel: "noopener" }, text);
}
