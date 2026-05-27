import { revalidatePath } from "next/cache";
import { NextRequest, NextResponse } from "next/server";
import { getPayload } from "payload";
import config from "@payload-config";
import { SUPPORTED_LANGS, SupportedLang } from "@/shared/constants/supported_langs";

export async function POST(req: NextRequest) {
  const secret = req.headers.get("x-revalidate-secret");
  if (!secret || secret !== process.env.PAYLOAD_SECRET) {
    return NextResponse.json({ ok: false }, { status: 401 });
  }

  const payload = await getPayload({ config });

  for (const lang of SUPPORTED_LANGS) {
    // Homepage
    revalidatePath(`/${lang}`);

    // All slug pages
    const result = await payload.find({
      collection: "pages",
      locale: lang as SupportedLang,
      draft: false,
      limit: 1000,
      select: { slug: true },
    });

    for (const page of result.docs) {
      if (page.slug) {
        revalidatePath(`/${lang}/${page.slug}`);
      }
    }
  }

  return NextResponse.json({ ok: true });
}
