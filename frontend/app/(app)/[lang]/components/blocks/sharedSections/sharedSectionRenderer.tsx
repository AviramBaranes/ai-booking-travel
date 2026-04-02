import { SharedSection, SharedSectionRefBlock } from "@/payload-types";
import { TypedSection } from "@/shared/types/payload";
import { NewsletterSection } from "./NewsletterSection";
import { StatsSection } from "./StatsSection";
import { SuppliersSection } from "./SuppliersSection";

interface SharedSectionRendererProps {
  section: SharedSection;
}
export function SharedSectionRenderer({ section }: SharedSectionRendererProps) {
  switch (section.type) {
    case "newsletter":
      return (
        <NewsletterSection section={section as TypedSection<"newsletter">} />
      );
    case "stats":
      return <StatsSection section={section as TypedSection<"stats">} />;
    case "suppliers":
      return (
        <SuppliersSection section={section as TypedSection<"suppliers">} />
      );
    default:
      return <div>[Unknown SharedSection: {section.type}]</div>;
  }
}
