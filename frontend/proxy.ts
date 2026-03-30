import { NextRequest, NextResponse } from "next/server";

const SUPPORTED_LANGS = ["he", "en"];

export function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const lang = pathname.split("/")[1];
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
