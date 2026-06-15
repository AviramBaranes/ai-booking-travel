import Image from "next/image";

export function PagesDecorations() {
  return (
    <div
      aria-hidden="true"
      className="pointer-events-none absolute inset-0 -z-10 overflow-hidden"
    >
      <div className="absolute top-100 right-0">
        <Image
          src="/assets/pages/road.png"
          alt=""
          width={225}
          height={100}
        />
      </div>

      <div className="absolute top-400 -left-25">
        <Image
          src="/assets/pages/orange-decor.png"
          alt=""
          width={500}
          height={100}
        />
      </div>
    </div>
  );
}