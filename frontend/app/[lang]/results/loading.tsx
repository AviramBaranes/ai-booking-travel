import Image from "next/image";

const brands = [
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

export default function LoadingPage() {
  return (
    <div className="text-center">
      <p className="mt-40">
        דמי ביטול 5% או 100 ש״ח הנמוך מביניהם עד 48 שעות לפני איסוף הרכב
      </p>
      <hr className="w-1/12 mx-auto bg-[#336CAE] h-1 rounded border-none mt-3" />

      <div className="mt-5 self-stretch text-center justify-start text-indigo-950 text-5xl font-black leading-12">
        מחפשים לכם בכל מקום בעולם
      </div>
      <div
        dir="ltr"
        className="relative flex items-center py-30 overflow-hidden"
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

        <div className="absolute z-10 bottom-13 animate-slide-right pointer-events-none">
          <div
            id="blur-bg"
            className="absolute rounded-full backdrop-blur-[1.5px] w-39 h-39 top-22 left-22 -translate-x-1/2 -translate-y-1/2"
          ></div>
          <Image
            src="/assets/magnifying-glass.svg"
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
