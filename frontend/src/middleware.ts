import { NextResponse } from "next/server";

import type { NextRequest } from "next/server";

// make sure they have some key.
// real authentication will happen within the api routes
function authenticate(refresh_token: string | undefined): boolean {
  // Description: shallow authentication to ensure a user has a refresh token
  // args: refresh_token -- users refresh token
  // returns: boolean : false if lacking refresh token, true otherwise

  if (refresh_token === undefined || refresh_token === "") {
    return false;
  }
  return true;
}

export function middleware(request: NextRequest): NextResponse | undefined {
  // Description:
  // args: request: incoming request
  // returns: NextResponse | undefined : either a redirect or allows the response through, unmodified

  // verifys if you have a token
  const authenticated = authenticate(
    request.cookies.get("refresh_token")?.value
  );

  if (!authenticated) {
    // push user back to login page
    return NextResponse.redirect(new URL("/login", request.url));
  }
}

export const config = {
  // matches requests to "/"
  matcher: "/",
};
