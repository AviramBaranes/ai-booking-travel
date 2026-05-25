import { NextRequest, NextResponse } from "next/server";

const SUPPORTED_LANGS = ["he", "en"];

export function proxy(request: NextRequest) {
  const { pathname, searchParams } = request.nextUrl;
  const lang = pathname.split("/")[1];

  if (searchParams.has("payload_preview")) {
    const url = request.nextUrl.clone();
    url.searchParams.delete("payload_preview");
    const res = NextResponse.redirect(url);
    res.cookies.set("payload-preview", "1", { maxAge: 60 * 60, path: "/" });
    return res;
  }

  const response = NextResponse.next();

  if (lang && SUPPORTED_LANGS.includes(lang)) {
    response.headers.set("X-Lang", lang);
  }

  if (request.nextUrl.pathname === "/cms/login") {
    const loginPage = new URL("/he", request.url);
    return NextResponse.redirect(loginPage);
  }

  return response;
}

export const config = {
  matcher: ["/((?!_next|api|favicon.ico).*)"],
};
