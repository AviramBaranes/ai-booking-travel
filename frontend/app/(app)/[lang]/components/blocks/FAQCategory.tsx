import type { FAQBlock } from "@/payload-types";
import { FAQItem } from "./FAQItem";

type Category = NonNullable<FAQBlock["categories"]>[number];

interface FAQCategoryProps {
  category: Category;
}

export function FAQCategory({ category }: FAQCategoryProps) {
  return (
    <div className="flex flex-col gap-4">
      {category.heading && (
        <h3 className="text-2xl font-bold leading-12 text-navy">
          {category.heading}
        </h3>
      )}
      <div className="flex flex-col">
        {category.items?.map((item, index) => (
          <FAQItem key={item.id ?? index} item={item} />
        ))}
      </div>
    </div>
  );
}
