import { NextRequest, NextResponse } from "next/server";

const SUPPORTED_LANGS = ["he", "en"];

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const lang = pathname.split("/")[1];
  const response = NextResponse.next();

  if (lang && SUPPORTED_LANGS.includes(lang)) {
    response.headers.set("X-Lang", lang);
  }

  return response;
}

export const config = {
  matcher: ["/((?!_next|api|favicon.ico).*)"],
};
