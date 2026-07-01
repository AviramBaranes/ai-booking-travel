import { FAQBlock as FAQBlockData } from "@/payload-types";
import { SectionHeader } from "../SectionHeader";
import { cn } from "@/lib/utils";
import { FAQCategory } from "./FAQCategory";

interface FAQBlockProps {
  data: FAQBlockData;
  className?: string;
}

export function FAQBlock({ data, className }: FAQBlockProps) {
  return (
    <div
      className={cn(
        "mx-5 mb-20 flex flex-col items-center gap-12 max-lg:w-auto! lg:mx-auto",
        className,
      )}
      style={{ width: data.width ?? "100%" }}
    >
      <SectionHeader
        pillText={data.eyebrow ?? ""}
        title={data.title ?? ""}
        subtitle={data.subtitle}
      />
      <div className="flex w-full flex-col gap-12">
        {data.categories?.map((category, index) => (
          <FAQCategory
            key={category.id ?? index}
            category={category}
            columns={data.columns ?? 1}
          />
        ))}
      </div>
    </div>
  );
}
