import Image from "next/image";

interface CarChecksProps {
  checks: string[];
}
export function CarChecks({ checks }: CarChecksProps) {
  return (
    <div className="flex flex-col mt-4">
      {checks.map((check) => (
        <div className="flex gap-2 items-center">
          <Image
            src="/assets/icons/V.svg"
            alt="Checked Icon"
            width={16}
            height={4}
            className="w-4"
          />
          <span className="type-paragraph text-navy">{check}</span>
        </div>
      ))}
    </div>
  );
}
