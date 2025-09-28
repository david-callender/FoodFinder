import { NextResponse } from "next/server";

import type { NextRequest } from "next/server";

// make sure they have some key.
// real authentication will happen within the api routes
function authenticate(access_token: string | undefined): boolean {
  if (access_token === undefined || access_token === "") {
    return false;
  }
  return true;
}

export function middleware(request: NextRequest): NextResponse | undefined {
  // verifys if you have a token
  const authenticated = authenticate(
    request.cookies.get("refresh_token")?.value
  );

  if (!authenticated) {
    return NextResponse.redirect(new URL("/login", request.url));
  }
}

export const config = {
  // matches requests to "/"
  matcher: "/",
};
