import Image from "next/image";

interface CarChecksProps {
  checks: { text: string; image: string }[];
}
export function CarChecks({ checks }: CarChecksProps) {
  return (
    <div className="flex flex-col mt-4">
      {checks.map((check) => (
        <div className="flex gap-2 items-center" key={check.text}>
          <Image
            src={check.image}
            alt={check.text}
            width={24}
            height={24}
            className="w-6 h-6"
          />
          <span className="type-paragraph text-navy">{check.text}</span>
        </div>
      ))}
    </div>
  );
}
