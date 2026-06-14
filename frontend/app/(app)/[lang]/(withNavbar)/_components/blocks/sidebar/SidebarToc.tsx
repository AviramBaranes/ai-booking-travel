"use client";

import { slugify } from "@/shared/utils/slugify";
import { useState, useEffect, useRef } from "react";

type SidebarTocProps = {
  sections: { title: string; id?: string | null }[];
};

export function SidebarToc({ sections }: SidebarTocProps) {
  const isScrollingRef = useRef(true); // start locked to ignore observer until hash check

  // Scroll to hash on mount
  useEffect(() => {
    const hash = decodeURIComponent(window.location.hash.slice(1));

    if (hash) {
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

  const handleClick = (e: React.MouseEvent, anchor: string) => {
    e.preventDefault();
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
    <nav className="rounded-lg border-4 border-t-0 border-b-0 my-12 border-l-brand border-r-brand bg-background px-10 w-full">
      <h5 className="type-h5 text-navy mt-5">תוכן עניינים:</h5>
      <hr className="my-5" />
      <ol className="list-decimal ps-5">
        {sections.map((section) => {
          const anchor = slugify(section.title);
          return (
            <li key={section.id} className="mb-3 text-brand-blue">
              <a href={`#${anchor}`} onClick={(e) => handleClick(e, anchor)}>
                <span />
                {section.title}
              </a>
            </li>
          );
        })}
      </ol>
    </nav>
  );
}
