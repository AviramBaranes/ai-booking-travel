import Image from "next/image";

export function PagesDecorations() {
  return (
    <>
      <div className="absolute top-100">
        <Image
          src="/assets/pages/road.png"
          alt="Road"
          width={225}
          height={100}
        />
      </div>
      <div className="absolute top-400 -z-10 -left-25">
        <Image
          src="/assets/pages/orange-decor.png"
          alt="Orange Decoration"
          width={500}
          height={100}
        />
      </div>
    </>
  );
}
