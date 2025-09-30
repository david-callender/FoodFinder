// jose runs in the "edge runtime"
// which basically means it runs closer to the browser level
// this is needed in middleware because middleware is working
// more "barebones" so to speak with requests, etc. BEFORE any rendering
import { jwtVerify } from "jose";
import { NextResponse } from "next/server";

import type { NextRequest } from "next/server";

async function authenticate(
  refresh_token: string | undefined
): Promise<boolean> {
  // Description: shallow authentication to ensure a user has a refresh token
  // args: refresh_token -- users refresh token
  // returns: boolean : false if lacking refresh token, true otherwise

  if (refresh_token === undefined || refresh_token === "") {
    return false;
  }
  // need to encode secret for verification
  // Per chatgpt: jose is built on WEB Crypto API, the only library available for edge runtime.
  // this library requires secret keys to be encoded
  const refresh_secret = new TextEncoder().encode(
    process.env.refresh_key?.trim()
  );

  try {
    const { payload } = await jwtVerify(refresh_token, refresh_secret);
    console.log(payload);
    return true;
  } catch (error) {
    console.log(error);
    return false;
  }
}

export async function middleware(request: NextRequest): Promise<NextResponse> {
  // Description:
  // args: request: incoming request
  // returns: NextResponse | undefined : either a redirect or allows the response through, unmodified

  // verifys if you have a token
  const authenticated = await authenticate(
    request.cookies.get("refresh_token")?.value
  );

  if (!authenticated) {
    // push user back to login page
    return NextResponse.redirect(new URL("/login", request.url));
  }
  // helps with a cleaner function signature
  // functionally no difference with or without return statement
  return NextResponse.next();
}

export const config = {
  // matches requests to "/"
  matcher: "/",
};
