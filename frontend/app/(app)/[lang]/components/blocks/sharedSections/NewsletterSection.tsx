import Image from "next/image";
import Link from "next/link";
import { TypedSection, Populated } from "@/shared/types/payload";

type NewsletterSectionProps = {
  section: TypedSection<"newsletter">;
};

export function NewsletterSection({ section }: NewsletterSectionProps) {
  const {
    title,
    subtitle,
    benefits,
    formTitle,
    formSubTitle,
    emailPlaceholder,
    submitButtonLabel,
    consentLabel,
    privacyTextBeforeLink,
    privacyLinkLabel,
    privacyPage,
  } = section.newsletter;

  const privacySlug = (privacyPage as Populated<typeof privacyPage>)?.slug;

  return (
    <section
      className="relative w-2/3 mx-auto overflow-hidden rounded-[20px] px-7 py-11 bg-cover bg-center bg-no-repeat"
      style={{ backgroundImage: "url('/assets/newsletter/newsletter-bg.png')" }}
    >
      <div className="relative flex flex-col items-center gap-10 md:flex-row md:items-center md:justify-center md:gap-20">
        {/* ── Right (in RTL): hero text + benefits ── */}
        <div className="flex w-full flex-col items-start gap-5 text-start md:w-100.25">
          <h2 className="text-4xl font-black leading-tight text-background">
            {title}
          </h2>
          {subtitle && <p className="text-base text-gray-400">{subtitle}</p>}
          {benefits && benefits.length > 0 && (
            <div className="flex flex-wrap items-center gap-5">
              {benefits.map((b) => (
                <div key={b.id} className="flex items-start gap-1.5">
                  <span className="text-xs font-bold text-green-600">✓</span>
                  <span className="text-sm text-white">{b.text}</span>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* ── Left (in RTL): form ── */}
        <form className="flex flex-1 flex-col items-start gap-5">
          {/* Heading */}
          <div className="flex flex-col items-start gap-1.5 text-start">
            {formTitle && (
              <h3 className="text-xl font-bold text-white">{formTitle}</h3>
            )}
            {formSubTitle && (
              <p className="text-sm text-gray-400">{formSubTitle}</p>
            )}
          </div>

          {/* Email input + submit */}
          <div className="flex w-full items-center gap-2">
            <div className="relative flex-1">
              <Image
                src="/assets/newsletter/mail-icon.svg"
                alt=""
                width={16}
                height={16}
                aria-hidden
                className="pointer-events-none absolute top-1/2 inset-s-4 -translate-y-1/2"
              />
              <input
                type="email"
                placeholder={emailPlaceholder ?? "כתובת הדואר האלקטרוני שלכם"}
                className="w-full rounded-[10px] border border-border-light bg-white px-4.5 py-4 ps-10 text-sm text-foreground shadow-[0px_2px_8px_0px_rgba(0,0,0,0.06)] placeholder:text-gray-400"
              />
            </div>
            <button
              type="submit"
              className="h-13 w-30 shrink-0 cursor-pointer rounded-[10px] bg-brand text-[15px] font-bold text-white shadow-[0px_2px_8px_0px_rgba(0,0,0,0.06)] transition-opacity hover:opacity-90"
            >
              {submitButtonLabel}
            </button>
          </div>

          {/* Consent + privacy */}
          <div className="flex w-full flex-col items-start gap-2">
            {consentLabel && (
              <label className="flex items-center gap-2 text-sm text-white cursor-pointer">
                <input
                  type="checkbox"
                  className="size-4.5 cursor-pointer rounded border-[1.5px] border-border-light bg-white accent-brand"
                />
                <span>{consentLabel}</span>
              </label>
            )}
            {privacyTextBeforeLink && (
              <p className="text-sm text-gray-400">
                {privacyTextBeforeLink}
                {privacyLinkLabel && privacySlug && (
                  <Link
                    target="_blank"
                    href={privacySlug}
                    className="cursor-pointer font-bold text-brand-blue underline"
                  >
                    {" "}
                    {privacyLinkLabel}
                  </Link>
                )}
              </p>
            )}
          </div>
        </form>
      </div>
    </section>
  );
}
