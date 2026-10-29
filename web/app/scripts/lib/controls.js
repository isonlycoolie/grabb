import { h, icon } from "./dom.js";

/**
 * choices renders a one-of-many group. Each option is
 * [value, label, iconName?, disabled?]; the group keeps its own state.
 */
function choices(className, options, current, onChange, label) {
  const buttons = options.map(([value, text, iconName, disabled]) =>
    h(
      "button",
      { type: "button", "aria-pressed": String(value === current), disabled, onclick: () => select(value) },
      iconName && icon(iconName, 14),
      text,
    ),
  );
  function select(value) {
    buttons.forEach((button, i) => button.setAttribute("aria-pressed", String(options[i][0] === value)));
    onChange(value);
  }
  return h("div", { class: className, role: "group", "aria-label": label }, buttons);
}

/** segmented is a compact switch for two to four options. */
export const segmented = (options, current, onChange, label) =>
  choices("segmented", options, current, onChange, label);

/** chips is a wrapping row of options, for quality and format. */
export const chips = (options, current, onChange, label) => choices("chips", options, current, onChange, label);

let region;

/** toast shows a short message in the corner for four seconds. */
export function toast(message, kind = "info") {
  region ??= document.body.appendChild(h("div", { class: "toasts", role: "status", "aria-live": "polite" }));
  const el = h("div", { class: "toast", "data-kind": kind }, icon(kind === "error" ? "alert" : "check"), h("span", {}, message));
  region.append(el);
  setTimeout(() => {
    el.classList.add("leaving");
    el.addEventListener("transitionend", () => el.remove(), { once: true });
  }, 4000);
}
