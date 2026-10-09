"use client";

import { useTranslations } from "next-intl";
import Image from "next/image";
import { Fragment, useEffect, useState } from "react";

const brands = [
  "Hertz.png",
  "Alamo.svg",
  "Europcar.svg",
  "Budget.svg",
  "Avis.svg",
  "Thrifty.svg",
  "Enterprise.svg",
  "Dollar.svg",
  "National.png",
  "Maggiore.png",
  "Green-Motion.png",
];

export function ResultsLoading() {
  return (
    <div className="text-center pt-40 pb-50">
      {/* <p className="mt-40">{t("loadingMessage")}</p>
      <hr className="w-1/12 mx-auto bg-[#336CAE] h-1 rounded border-none mt-3" /> */}

      <LoadingPageTitle />
      <div
        dir="ltr"
        className="relative flex items-center py-20 lg:py-30 overflow-hidden"
      >
        <div className="w-full overflow-hidden">
          <div className="animate-scroll-left">
            {[...brands, ...brands].map((logo, i) => (
              <div
                key={i}
                className="shrink-0 mx-4 flex items-center justify-center bg-white rounded-xl shadow-md p-5 w-44 h-24"
              >
                <Image
                  src={`/suppliers/${logo}`}
                  alt={logo.replace(/\.(svg|png)$/, "")}
                  width={120}
                  height={60}
                  className="object-contain max-h-16 w-auto h-auto"
                />
              </div>
            ))}
          </div>
        </div>

        <div className="absolute z-10 bottom-3 lg:bottom-13 animate-slide-right pointer-events-none">
          <div
            id="blur-bg"
            className="absolute rounded-full backdrop-blur-[1.5px] w-39 h-39 top-22 left-22 -translate-x-1/2 -translate-y-1/2"
          ></div>
          <Image
            src="/assets/loader/magnifying-glass.svg"
            alt="Magnifying glass"
            width={200}
            height={200}
            loading="eager"
          />
        </div>
      </div>
    </div>
  );
}

const WORD_ANIMATION_MS = 600;
const WORD_STAGGER_MS = 80;
const TITLE_HOLD_MS = 2800;

function LoadingPageTitle() {
  const t = useTranslations("ResultsPage");
  const titles = t.raw("loadingTitles") as string[];
  const [index, setIndex] = useState(0);
  const [leaving, setLeaving] = useState(false);

  const title = titles[index];
  const words = title.split(" ");
  // Time until the last word has finished entering (or leaving).
  const sequenceMs = WORD_ANIMATION_MS + (words.length - 1) * WORD_STAGGER_MS;

  // Words rise in one by one, the title holds, then the words lift out one by
  // one and the next title takes over.
  useEffect(() => {
    const id = leaving
      ? setTimeout(() => {
          setIndex((i) => (i + 1) % titles.length);
          setLeaving(false);
        }, sequenceMs)
      : setTimeout(() => setLeaving(true), sequenceMs + TITLE_HOLD_MS);
    return () => clearTimeout(id);
  }, [leaving, sequenceMs, titles.length]);

  return (
    <div className="mt-5 self-stretch text-center justify-start text-indigo-950 text-3xl mx-5 lg:text-5xl font-black leading-12 min-h-24 lg:min-h-12">
      <span className="sr-only">{title}</span>
      <span key={index} aria-hidden>
        {words.map((word, i) => (
          <Fragment key={i}>
            {i > 0 && " "}
            <span
              className={leaving ? "animate-word-out" : "animate-word-in"}
              style={{
                animationDuration: `${WORD_ANIMATION_MS}ms`,
                animationDelay: `${i * WORD_STAGGER_MS}ms`,
              }}
            >
              {word}
            </span>
          </Fragment>
        ))}
      </span>
    </div>
  );
}
