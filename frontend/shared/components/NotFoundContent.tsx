import Link from "next/link";

interface NotFoundProps {
  title: string;
  subtitle: string;
  buttonText: string;
  homepageUrl: string;
}

export function NotFoundContent({
  title,
  subtitle,
  buttonText,
  homepageUrl,
}: NotFoundProps) {
  return (
    <div className="flex mt-30 mb-10 items-center justify-center bg-background">
      <div className="flex flex-col items-center gap-6.5 w-175 max-w-full px-4">
        <h1 className="text-[220px] leading-none font-black text-navy text-center">
          404
        </h1>
        <h3 className="type-h3 text-navy text-center">{title}</h3>
        <h6 className="type-h6 text-muted text-center">{subtitle}</h6>
        <div className="bg-brand h-0.75 w-20 rounded-sm" />
        <Link
          href={homepageUrl}
          className="bg-brand text-white font-bold type-paragraph px-5 py-3.5 rounded-[10px] shadow-subtle whitespace-nowrap"
        >
          {buttonText}
        </Link>
      </div>
    </div>
  );
}
