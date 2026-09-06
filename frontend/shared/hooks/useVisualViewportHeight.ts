import { useEffect, useState } from "react";

/**
 * Height in pixels of the area the user can actually see.
 *
 * iOS Safari's viewport units all measure past the bottom toolbar. Measured on
 * iOS 18.7 with a 548px visible area: 100svh = 626.5px, 100dvh = 627px,
 * 100lvh = 717.5px. Anything sized with those units puts its lower edge behind
 * the toolbar, which is how full-height sheets lose their footer. visualViewport
 * reports the real box, and tracks pinch-zoom and the on-screen keyboard too.
 *
 * Undefined until measured, and on browsers without visualViewport — callers
 * should keep a CSS fallback for that case.
 */
export function useVisualViewportHeight() {
  const [height, setHeight] = useState<number>();

  useEffect(() => {
    const viewport = window.visualViewport;
    if (!viewport) return;

    const update = () => setHeight(viewport.height);
    update();
    viewport.addEventListener("resize", update);
    return () => viewport.removeEventListener("resize", update);
  }, []);

  return height;
}
