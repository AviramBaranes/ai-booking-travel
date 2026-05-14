import { SharedSectionRefBlock } from "@/payload-types";
import { PropsWithChildren } from "react";

const spacingTopClass: Record<
  NonNullable<NonNullable<SharedSectionRefBlock["overrides"]>["spacingTop"]>,
  string
> = {
  default: "mt-16",
  none: "mt-0",
  sm: "mt-8",
  md: "mt-16",
  lg: "mt-24",
};

const spacingBottomClass: Record<
  NonNullable<NonNullable<SharedSectionRefBlock["overrides"]>["spacingBottom"]>,
  string
> = {
  default: "mb-16",
  none: "mb-0",
  sm: "mb-8",
  md: "mb-16",
  lg: "mb-24",
};

export function SharedSectionWrapper({
  children,
  overrides,
}: PropsWithChildren<{ overrides: SharedSectionRefBlock["overrides"] }>) {
  const mt = spacingTopClass[overrides?.spacingTop ?? "md"];
  const mb = spacingBottomClass[overrides?.spacingBottom ?? "md"];

  return <div className={`${mt} ${mb}`}>{children}</div>;
}
