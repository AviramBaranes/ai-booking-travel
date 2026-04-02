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
        <p className="text-[220px] leading-none font-black text-navy text-center">
          404
        </p>
        <h1 className="text-[36px] leading-normal font-bold text-navy text-center">
          {title}
        </h1>
        <p className="text-[18px] leading-normal text-[#a0a3b8] text-center">
          {subtitle}
        </p>
        <div className="bg-brand h-0.75 w-20 rounded-sm" />
        <Link
          href={homepageUrl}
          className="bg-navy text-white font-bold text-[15px] px-32.5 py-3.5 rounded-[10px] shadow-[0px_2px_8px_0px_rgba(0,0,0,0.06)] whitespace-nowrap"
        >
          {buttonText}
        </Link>
      </div>
    </div>
  );
}
