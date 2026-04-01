import { SharedSection } from "@/payload-types";
import { TypedSection } from "@/shared/types/payload";
import { NewsletterSection } from "./NewsletterSection";

export function SharedSectionRenderer({ section }: { section: SharedSection }) {
  switch (section.type) {
    case "newsletter":
      return (
        <NewsletterSection section={section as TypedSection<"newsletter">} />
      );
    default:
      return <div>[Unknown SharedSection: {section.type}]</div>;
  }
}
