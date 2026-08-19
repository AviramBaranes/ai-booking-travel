import { revalidatePath } from "next/cache";
import { NextRequest, NextResponse } from "next/server";
import { SUPPORTED_LANGS } from "@/shared/constants/supported_langs";
import type { RevalidateRequest } from "@/CMS/hooks/revalidate";

export async function POST(req: NextRequest) {
  const secret = req.headers.get("x-revalidate-secret");

  if (!secret || secret !== process.env.PAYLOAD_SECRET) {
    return NextResponse.json({ ok: false }, { status: 401 });
  }

  const body: RevalidateRequest | null = await req.json().catch(() => null);

  revalidatePath("/", "layout");
  revalidatePath("/[lang]", "page");
  revalidatePath("/[lang]/[slug]", "page");
  // Content changes alter the URL set and lastModified dates.
  revalidatePath("/sitemap.xml");

  if (body?.collection === "pages" && body?.slug) {
    for (const lang of SUPPORTED_LANGS) {
      revalidatePath(`/${lang}/${body.slug}`);
    }
  }

  if (body?.collection === "blog-posts") {
    // Route-pattern form covers every instance, which sidesteps the
    // encoded/decoded ambiguity of Hebrew slugs in a literal path.
    revalidatePath("/[lang]/blog/[slug]", "page");
    revalidatePath("/[lang]/blog/page/[page]", "page");
  }

  return NextResponse.json({ ok: true });
}
