// Applies the saved theme before first paint. Loaded as a classic script so
// it runs before the stylesheets render; the settings page calls
// window.grabbTheme to change it.
(function () {
  var key = "grabb-theme";
  var media = window.matchMedia("(prefers-color-scheme: dark)");
  var saved = null;
  try {
    saved = localStorage.getItem(key);
  } catch (e) {}

  function apply() {
    var dark = saved === "dark" || (saved !== "light" && media.matches);
    document.documentElement.dataset.theme = dark ? "dark" : "light";
  }

  apply();
  media.addEventListener("change", apply);

  window.grabbTheme = {
    get: function () {
      return saved || "system";
    },
    set: function (theme) {
      saved = theme === "system" ? null : theme;
      try {
        if (saved) localStorage.setItem(key, saved);
        else localStorage.removeItem(key);
      } catch (e) {}
      apply();
    },
  };
})();
