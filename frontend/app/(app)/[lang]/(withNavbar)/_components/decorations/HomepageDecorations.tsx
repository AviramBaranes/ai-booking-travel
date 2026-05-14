import Image from "next/image";

export function HomepageDecorations() {
  return (
    <div dir="rtl">
      <div className="absolute -z-10 top-150">
        <Image
          src="/assets/home/paper-plane.png"
          alt="Paper Plane"
          width={360}
          height={360}
          className="w-90 h-90"
        />
      </div>
      <div className="absolute top-325 -z-10 -left-10">
        <Image
          src="/assets/home/suitcase.png"
          alt="Suitcase"
          width={380}
          height={360}
          className="w-95 h-90"
        />
      </div>
      <div className="absolute -z-10 top-550">
        <Image
          src="/assets/home/car.png"
          alt="Car"
          width={225}
          height={225}
          className="w-56 h-56"
        />
      </div>
      <div className="absolute -z-10 top-550">
        <Image
          src="/assets/home/road.png"
          alt="Road"
          width={225}
          height={225}
        />
      </div>
      <div className="absolute top-430 -z-10 -left-25">
        <Image
          src="/assets/home/airplane.png"
          alt="Airplane"
          width={500}
          height={500}
          className="w-125 h-125"
        />
      </div>
      <div className="absolute top-740 -z-10 -left-10">
        <Image
          src="/assets/home/location.png"
          alt="Location"
          width={300}
          height={300}
          className="w-75 h-75"
        />
      </div>
      {/* <div className="absolute top-800 -z-10 -left-25">
        <Image
          src="/assets/pages/orange-decor.png"
          alt="Orange Decoration"
          width={500}
          height={100}
        />
      </div> */}
    </div>
  );
}
