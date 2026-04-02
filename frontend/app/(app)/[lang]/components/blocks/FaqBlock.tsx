import { FAQBlock as FAQBlockData } from "@/payload-types";
import { SectionHeader } from "./SectionHeader";
import { FAQCategory } from "./FAQCategory";

interface FAQBlockProps {
  data: FAQBlockData;
}

export function FAQBlock({ data }: FAQBlockProps) {
  return (
    <div className="flex flex-col items-center gap-12 w-1/2 mx-auto mb-20">
      <SectionHeader
        pillText={data.eyebrow ?? ""}
        title={data.title ?? ""}
        subtitle={data.subtitle}
      />
      <div className="flex w-full max-w-250 flex-col gap-12">
        {data.categories?.map((category, index) => (
          <FAQCategory key={category.id ?? index} category={category} />
        ))}
      </div>
    </div>
  );
}
