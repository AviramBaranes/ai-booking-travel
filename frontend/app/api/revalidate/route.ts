import { revalidatePath } from "next/cache";
import { NextRequest, NextResponse } from "next/server";

export async function POST(req: NextRequest) {
  const secret = req.headers.get("x-revalidate-secret");
  if (!secret || secret !== process.env.PAYLOAD_SECRET) {
    return NextResponse.json({ ok: false }, { status: 401 });
  }

  revalidatePath("/", "layout");

  return NextResponse.json({ ok: true });
}
