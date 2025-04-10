import { NextResponse } from "next/server";

export function middleware(req) {
  const { pathname } = req.nextUrl;

  if (pathname.startsWith("/admin")) {
    const idToken = req.cookies.get("id_token");
    if (!idToken) {
      const loginUrl = new URL("http://localhost:8080/login");
      loginUrl.searchParams.set("redirect_uri", req.nextUrl.href);
      return NextResponse.redirect(loginUrl);
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/admin/:path*"],
};
