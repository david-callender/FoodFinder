import { useState } from "react";

import type { ChangeEvent, FC } from "react";

function wrapPhoneNumber(phoneNumber: string): string {
  if (phoneNumber.length !== 10) {
    return "";
  }

  const areaCode = phoneNumber.slice(0, 3);
  const officeCode = phoneNumber.slice(3, 6);
  const lineNumber = phoneNumber.slice(6, 10);

  return `(${areaCode})-${officeCode}-${lineNumber}`;
}

export const PhoneNumberInput: FC = () => {
  // setting state
  // phone number is stored purely as a string of integers
  const [phoneNumber, setPhoneNumber] = useState("");

  // removes illicit characters from phoneNumber
  function handlePhoneNumber(event: ChangeEvent<HTMLInputElement>): void {
    // matches any non-digit character
    const blacklistRegex = /[^0-9]/g;
    const rawPhoneNumber = event.target.value;

    const strippedPhoneNumber = rawPhoneNumber.replaceAll(blacklistRegex, "");
    // matches numbers with > 10 digits
    const limitedPhoneNumber = strippedPhoneNumber.replace(
      /\d{11,}/,
      strippedPhoneNumber.slice(0, 10)
    );

    setPhoneNumber(limitedPhoneNumber);
  }

  // TODO : switching components returns the user to end the string if they try to delete from
  // the middle of the masked phone number. Is there a fix for this? Is this an actual issue?
  return phoneNumber.length === 10 ? (
    <input
      type="tel"
      value={wrapPhoneNumber(phoneNumber)}
      className="m-3 place-self-center rounded-lg border-4 border-gray-200 bg-gray-200 p-0.5 text-black focus:border-gray-400"
      name="phone-number"
      placeholder="XXX-XXX-XXXX"
      pattern="[0-9]{3}[0-9]{3}[0-9]{4}"
      onChange={(event) => {
        handlePhoneNumber(event);
      }}
    />
  ) : (
    <input
      type="tel"
      value={phoneNumber}
      className="m-3 place-self-center rounded-lg border-4 border-gray-200 bg-gray-200 p-0.5 text-black"
      name="phone-number"
      placeholder="XXX-XXX-XXXX"
      pattern="[0-9]{3}[0-9]{3}[0-9]{4}"
      onChange={(event) => {
        handlePhoneNumber(event);
      }}
    />
  );
};
