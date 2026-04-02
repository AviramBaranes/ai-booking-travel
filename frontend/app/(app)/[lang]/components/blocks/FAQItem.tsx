import type { FAQBlock } from "@/payload-types";
import { RichText } from "@payloadcms/richtext-lexical/react";
import { ChevronDown } from "lucide-react";

type Item = NonNullable<
  NonNullable<FAQBlock["categories"]>[number]["items"]
>[number];

interface FAQItemProps {
  item: Item;
}

export function FAQItem({ item }: FAQItemProps) {
  return (
    <details className="group border-b border-border-light py-4">
      <summary className="flex cursor-pointer list-none items-center justify-between gap-4 [&::-webkit-details-marker]:hidden">
        <span className="text-base font-medium text-navy">{item.question}</span>
        <ChevronDown
          size={20}
          className="shrink-0 text-navy transition-transform duration-200 group-open:rotate-180"
        />
      </summary>
      <div className="py-4 text-sm leading-relaxed">
        <RichText data={item.answer} />
      </div>
    </details>
  );
}
