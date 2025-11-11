"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

import { signup } from "@/db/signup";

import type { ChangeEvent, FC, FormEvent } from "react";

function wrapPhoneNumber(phoneNumber: string): string {
  // Purpose: to mask phone number in input component
  // Args:
  // phoneNumber: string - string representing a 10 digit phone number
  // Returns:
  // string - formatted phone number of (XXX) - XXX - XXXX for masking in an input component

  if (phoneNumber.length !== 10) {
    return "";
  }

  const areaCode = phoneNumber.slice(0, 3);
  const officeCode = phoneNumber.slice(3, 6);
  const lineNumber = phoneNumber.slice(6, 10);

  return `(${areaCode})-${officeCode}-${lineNumber}`;
}

function removeBlacklistCharacters(phoneNumber: string): string {
  // Purpose: to remove non digit and extraneous characters from an input phone number
  // Args:
  // phoneNumber: string - string representing a 10 digit phone number
  // Returns:
  // string - <= 10 character string representing a possible phone number

  // matches non-digit characters
  const blacklistRegex = /[^0-9]/g;

  const strippedPhoneNumber = phoneNumber.replaceAll(blacklistRegex, "");
  // matches numbers with > 10 digits
  const limitedPhoneNumber = strippedPhoneNumber.replace(
    /\d{11,}/,
    strippedPhoneNumber.slice(0, 10)
  );

  return limitedPhoneNumber;
}

export const SignUpBox: FC = () => {
  const router = useRouter();

  // for when we add display name
  //const [displayName, setDisplayName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [phoneNumber, setPhoneNumber] = useState("");
  const [displayName, setDisplayName] = useState("");

  function changePhoneNumber(e: ChangeEvent<HTMLInputElement>): void {
    // Purpose : controlling state for phone number field in form
    // Args:
    // event : ChangeEvent<HTMLInputElement> - event from Input element
    // Returns
    // void - changes phoneNumber state in place

    const rawPhoneNumber = e.target.value; // TODO [misc.] : further type/format verification?

    const limitedPhoneNumber = removeBlacklistCharacters(rawPhoneNumber);

    if (limitedPhoneNumber.length === 10) {
      setPhoneNumber(wrapPhoneNumber(limitedPhoneNumber));
    } else {
      setPhoneNumber(limitedPhoneNumber);
    }
  }

  async function handleSubmit(
    event: FormEvent<HTMLFormElement>
  ): Promise<void> {
    // Purpose : POSTing form/state values to server
    // Args:
    // event : ChangeEvent<HTMLInputElement> - event from Input element
    // Returns
    // void - posting data to server
    // prevents refresh of page
    event.preventDefault();

    await signup(email, password, displayName.trim());

    localStorage.setItem("displayName", displayName.trim());
    //  TODO [backend] : Handling a "user exists" error from backend (or if they already have cookies)
    // TODO [misc.] : Handle User phone numbers without a US country code

    router.push("/");
  }

  return (
    <form onSubmit={handleSubmit}>
      <div className="flex flex-col">
        <input
          type="text"
          value={displayName}
          onChange={(e) => {
            setDisplayName(e.target.value);
          }}
          className="m-3 place-self-center rounded-lg bg-gray-800 p-2"
          name="displayName"
          placeholder="Display Name"
          required
        />
        <input
          type="email"
          value={email}
          onChange={(e) => {
            setEmail(e.target.value);
          }}
          className="m-3 place-self-center rounded-lg bg-gray-800 p-2"
          name="email"
          placeholder="Email"
          required
        />
        <input
          type="password"
          value={password}
          onChange={(e) => {
            setPassword(e.target.value);
          }}
          className="m-3 place-self-center rounded-lg bg-gray-800 p-2"
          name="password"
          placeholder="Password"
          required
        />
        <input
          type="tel"
          value={phoneNumber}
          onChange={(e) => {
            changePhoneNumber(e);
          }}
          className="m-3 place-self-center rounded-lg bg-gray-800 p-2"
          name="phone-number"
          placeholder="XXX-XXX-XXXX"
          required
        />
        <button className="mx-auto rounded-xl bg-red-900 px-4 py-2 font-semibold shadow transition hover:cursor-pointer hover:bg-red-700">
          Sign Up
        </button>
        <p className="m-2 place-self-center text-xs">
          Already have an account?{" "}
          <a href="/login" className="text-blue-400">
            Login Here
          </a>
        </p>
      </div>
    </form>
  );
};
