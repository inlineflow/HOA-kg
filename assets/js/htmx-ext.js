"use strict";

document.querySelector("body").addEventListener("htmx:configRequest", (e) => {
  let s = `${e.detail.path}?`;
  if (e.detail.verb === "delete") {
    for (const p in e.detail.parameters) {
      let x = `${p}=${e.detail.parameters[p].join(",")}&`;
      s = s + x;
    }
  }
  e.detail.path = s;
});
