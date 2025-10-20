function overflowMenu(tree = document) {
  tree.querySelectorAll("[data-overflow-menu]").forEach((menuRoot) => {
    const button = menuRoot.querySelector("[aria-haspopup]");
    const menu = menuRoot.querySelector("[role=menu]");
    const items = [...menuRoot.querySelectorAll("[role=menuitem]")];

    const isOpen = () => !menu.hidden;
    items.forEach((item) => item.setAttribute("tabindex", -1));

    function toggleMenu(open = !isOpen()) {
      if (open) {
        menu.hidden = false;
        button.setAttribute("aria-expanded", "true");
        items[0].focus();
      } else {
        menu.hidden = true;
        button.setAttribute("aria-expanded", "false");
      }
    }

    toggleMenu(isOpen());
    button.addEventListener("click", (e) => {
      console.log(e.target);
      toggleMenu();
    });
    // menuRoot.addEventListener("blur", () => toggleMenu(false));

    window.addEventListener("click", function clickAway(event) {
      if (!menuRoot.isConnected) {
        window.removeEventListener("click", clickAway);
      }

      if (!menuRoot.contains(event.target)) {
        toggleMenu(false);
      }
    });

    const currentIndex = () => {
      const idx = items.indexOf(document.activeElement);
      if (idx === -1) {
        return 0;
      }
      return idx;
    };

    menu.addEventListener("keydown", (e) => {
      switch (e.key) {
        case "ArrowUp":
          items[currentIndex() - 1]?.focus();
          e.preventDefault();
          break;
        case "ArrowDown":
          items[currentIndex() + 1]?.focus();
          e.preventDefault();
          break;
        case "Space":
          items[currentIndex()].click();
          e.preventDefault();
          break;
        case "Home":
          items[0].focus();
          e.preventDefault();
          break;
        case "End":
          items[items.length - 1].focus();
          e.preventDefault();
          break;
        case "Escape":
          toggleMenu(false);
          button.focus();
      }
    });

    //
  });
}

addEventListener("htmx:load", (e) => overflowMenu(e.target));
