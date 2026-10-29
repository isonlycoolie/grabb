import { brands } from "./brands.js";
import { paths } from "./icons.js";

const svgNS = "http://www.w3.org/2000/svg";

/**
 * h builds an element. Props starting with "on" become listeners, the rest
 * become attributes; null and false props and children are skipped.
 *   h("button", { class: "btn", onclick }, icon("x"), "Cancel")
 */
export function h(tag, props = {}, ...children) {
  const el = document.createElement(tag);
  for (const [key, value] of Object.entries(props ?? {})) {
    if (value == null || value === false) continue;
    if (key.startsWith("on")) el.addEventListener(key.slice(2), value);
    else el.setAttribute(key, value === true ? "" : value);
  }
  el.append(...children.flat(Infinity).filter((c) => c != null && c !== false && c !== ""));
  return el;
}

/** icon returns one of the stroke icons from icons.js. */
export function icon(name, size = 16) {
  const svg = document.createElementNS(svgNS, "svg");
  const attrs = {
    viewBox: "0 0 24 24", width: size, height: size, fill: "none", stroke: "currentColor",
    "stroke-width": 1.75, "stroke-linecap": "round", "stroke-linejoin": "round", "aria-hidden": "true",
  };
  for (const [key, value] of Object.entries(attrs)) svg.setAttribute(key, value);
  svg.innerHTML = paths[name] ?? "";
  return svg;
}

/** thumbnail shows a remote image, falling back to an icon if it fails. */
export function thumbnail(src, fallback) {
  if (!src) return icon(fallback, 20);
  const img = h("img", { src, alt: "", loading: "lazy", referrerpolicy: "no-referrer" });
  img.addEventListener("error", () => img.replaceWith(icon(fallback, 20)), { once: true });
  return img;
}

/** siteMark shows a network's mark and name, from a yt-dlp extractor key. */
export function siteMark(site) {
  const brand = brands[String(site).toLowerCase()];
  if (!brand) return h("span", { class: "job-site" }, site);
  const svg = document.createElementNS(svgNS, "svg");
  for (const [key, value] of Object.entries({ viewBox: "0 0 24 24", width: 12, height: 12, fill: brand.color, "aria-hidden": "true" })) {
    svg.setAttribute(key, value);
  }
  svg.innerHTML = `<path d="${brand.path}"/>`;
  return h("span", { class: "job-site" }, svg, brand.name);
}

/** pageHeader is the title row at the top of a page. */
export function pageHeader(title, subtitle, ...actions) {
  return h(
    "header",
    { class: "page-header" },
    h("div", {}, h("h1", {}, title), subtitle && h("p", {}, subtitle)),
    actions.length > 0 && h("div", { class: "page-actions" }, actions),
  );
}
