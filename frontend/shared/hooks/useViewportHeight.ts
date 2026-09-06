import { useEffect, useState } from "react";

/**
 * Height in pixels of the area the user can actually see.
 *
 * iOS Safari's viewport units all measure past the bottom toolbar. Measured on
 * iOS 18.7 with a 548px visible area: 100svh = 626.5px, 100dvh = 627px,
 * 100lvh = 717.5px. Anything sized with those units puts its lower edge behind
 * the toolbar, which is how full-height sheets lose their footer.
 *
 * window.innerHeight reports the visible height correctly (548px in that same
 * measurement) and, unlike visualViewport, is unaffected by pinch/auto zoom and
 * by the on-screen keyboard — both of which would otherwise shrink a sheet to a
 * fraction of the screen.
 *
 * Undefined until measured, so callers should keep a CSS fallback.
 */
export function useViewportHeight() {
  const [height, setHeight] = useState<number>();

  useEffect(() => {
    const update = () => setHeight(window.innerHeight);
    update();
    window.addEventListener("resize", update);
    window.addEventListener("orientationchange", update);
    return () => {
      window.removeEventListener("resize", update);
      window.removeEventListener("orientationchange", update);
    };
  }, []);

  return height;
}
