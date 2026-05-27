import { revalidatePath } from "next/cache";
import { NextRequest, NextResponse } from "next/server";

export async function POST(req: NextRequest) {
  const secret = req.headers.get("x-revalidate-secret");
  if (!secret || secret !== process.env.PAYLOAD_SECRET) {
    return NextResponse.json({ ok: false }, { status: 401 });
  }

  // Use literal dynamic route patterns to bypass Hebrew URL encoding mismatches
  revalidatePath("/", "layout");
  revalidatePath("/[lang]", "page");
  revalidatePath("/[lang]/[slug]", "page");

  return NextResponse.json({ ok: true });
}
