import Image from "next/image";

interface InclusionsDisplayProps {
  title: string;
  inclusions: string[];
}

export function InclusionsDisplay({
  title,
  inclusions,
}: InclusionsDisplayProps) {
  return (
    <div className="min-h-115 bg-white border-cars-border rounded-xl shadow-card py-6">
      <h5 className="type-h5 text-navy mx-6">{title}</h5>
      <ul>
        {inclusions.map((inclusion) => (
          <li
            key={inclusion}
            className="type-paragraph text-navy mx-4 my-2 flex"
          >
            <Image
              src="/assets/icons/V.svg"
              alt="Included"
              width={28}
              height={28}
              className="inline w-7 h-7"
            />
            {inclusion}
          </li>
        ))}
      </ul>
    </div>
  );
}
