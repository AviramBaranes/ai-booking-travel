import type { FAQBlock } from "@/payload-types";
import { FAQItem } from "./FAQItem";

type Category = NonNullable<FAQBlock["categories"]>[number];

interface FAQCategoryProps {
  category: Category;
  columns: number;
}

export function FAQCategory({ category, columns }: FAQCategoryProps) {
  return (
    <div className="flex flex-col gap-4">
      {category.heading && (
        <h5 className="type-h5 text-navy">{category.heading}</h5>
      )}
      <div className="flex gap-x-8">
        {Array.from({ length: columns }, (_, col) => (
          <div key={col} className="flex flex-1 flex-col">
            {category.items
              ?.filter((_, i) => i % columns === col)
              .map((item, index) => (
                <FAQItem key={item.id ?? index} item={item} />
              ))}
          </div>
        ))}
      </div>
    </div>
  );
}
