import { revalidatePath } from "next/cache";
import { NextRequest, NextResponse } from "next/server";
import { SUPPORTED_LANGS } from "@/shared/constants/supported_langs";

export async function POST(req: NextRequest) {
  const secret = req.headers.get("x-revalidate-secret");

  if (!secret || secret !== process.env.PAYLOAD_SECRET) {
    return NextResponse.json({ ok: false }, { status: 401 });
  }

  const body = await req.json().catch(() => null);

  revalidatePath("/", "layout");
  revalidatePath("/[lang]", "page");
  revalidatePath("/[lang]/[slug]", "page");

  if (body?.collection === "pages" && body?.slug) {
    for (const lang of SUPPORTED_LANGS) {
      revalidatePath(`/${lang}/${body.slug}`);
    }
  }

  return NextResponse.json({ ok: true });
}