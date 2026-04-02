"use client";

import { slugify } from "@/shared/lang/slugify";
import { useState, useEffect, useRef } from "react";

type SidebarNavProps = {
  sections: { title: string; id?: string | null }[];
};

export function SidebarNav({ sections }: SidebarNavProps) {
  const [activeAnchor, setActiveAnchor] = useState<string | null>(
    slugify(sections[0]?.title) ?? null,
  );
  const isScrollingRef = useRef(true); // start locked to ignore observer until hash check

  // Scroll to hash on mount
  useEffect(() => {
    const hash = decodeURIComponent(window.location.hash.slice(1));

    if (hash) {
      setActiveAnchor(hash);

      // Push the scroll event to the end of the event loop to let Next.js finish hydrating
      setTimeout(() => {
        const el = document.getElementById(hash);
        if (el) {
          el.scrollIntoView({ behavior: "instant" });
        }
        // Unlock observer after the scroll settles
        isScrollingRef.current = false;
      }, 100); // 100ms is usually the sweet spot to beat Next.js scroll resets

      return;
    }

    isScrollingRef.current = false;
  }, []);

  useEffect(() => {
    const observers: IntersectionObserver[] = [];

    for (const section of sections) {
      const anchor = slugify(section.title);
      const el = document.getElementById(anchor);
      if (!el) continue;

      const observer = new IntersectionObserver(
        ([entry]) => {
          if (entry.isIntersecting && !isScrollingRef.current) {
            setActiveAnchor(anchor);
            history.replaceState(null, "", `#${anchor}`);
          }
        },
        { rootMargin: "-20% 0px -60% 0px" },
      );
      observer.observe(el);
      observers.push(observer);
    }

    return () => observers.forEach((o) => o.disconnect());
  }, [sections]);

  const handleClick = (e: React.MouseEvent, anchor: string) => {
    e.preventDefault();
    setActiveAnchor(anchor);
    isScrollingRef.current = true;

    const el = document.getElementById(anchor);
    el?.scrollIntoView({ behavior: "smooth" });
    history.replaceState(null, "", `#${anchor}`);

    // Re-enable observer only after scroll actually finishes
    const onScrollEnd = () => {
      requestAnimationFrame(() => {
        isScrollingRef.current = false;
      });
    };
    window.addEventListener("scrollend", onScrollEnd, { once: true });
  };

  return (
    <nav className="sticky top-40 w-72 shrink-0 rounded-xl border border-border-light bg-white p-6 shadow-sm">
      {sections.map((section) => {
        const anchor = slugify(section.title);
        const isActive = activeAnchor === anchor;
        return (
          <a
            key={section.id}
            href={`#${anchor}`}
            onClick={(e) => handleClick(e, anchor)}
            className={`flex items-center gap-2 rounded-lg p-4 text-lg font-semibold tracking-tight transition-colors ${
              isActive
                ? "border border-brand-blue bg-brand-blue/10 text-foreground"
                : "text-gray-600 hover:text-foreground"
            }`}
          >
            <span
              className={`h-1.5 w-1.5 shrink-0 rounded-full ${
                isActive ? "bg-brand-blue" : "bg-gray-300"
              }`}
            />
            {section.title}
          </a>
        );
      })}
    </nav>
  );
}
