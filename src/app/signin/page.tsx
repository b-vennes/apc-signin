"use client";

import Link from "next/link";
import Modal from "@/components/modal.tsx";
import { AgreementBlock } from "@/components/agreement-block.tsx";
import { FormEvent, useState } from "react";
import {
  AgreementCompleted,
  AgreementNeeded,
  AgreementState,
  AlreadyAgreed,
  useAgreement,
} from "@/lib/agreements.ts";
import { Debounce, useDebouncing } from "@/lib/debouncing.ts";
import { memberAgreed } from "@/app/actions.ts";

function getInputText(event: FormEvent<HTMLInputElement>): string {
  return (event.target as unknown as { value: string }).value;
}

function updateAgreementState(
  update: (value: AgreementState) => void,
  debouncing: Debounce,
  name: string,
  email: string,
) {
  async function checkPersonAgreed() {
    console.log(`Name: ${name}`);
    console.log(`Email: ${email}`);
    const agreed = await memberAgreed(name, email);

    console.log("Member agreed already? " + agreed);

    agreed ? update(AlreadyAgreed) : update(AgreementNeeded);
  }

  debouncing.reset(checkPersonAgreed);
}

export default function SignIn() {
  const [agreementOpen, setAgreementOpen] = useState(false);

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");

  // TODO implement agreement needed
  // this could all be one state with multiple pieces
  const [agreementState, setAgreementState] = useAgreement();

  const checkAgreeDebouncing = useDebouncing(500);

  const submitEnabled = name &&
    name !== "" &&
    email &&
    email !== "" &&
    (agreementState === AgreementCompleted || agreementState === AlreadyAgreed);

  return (
    <div
      id="sign-in-root"
      className="flex flex-col gap-2"
    >
      <div id="sign-in-header" className="bg-slate-200 p-4">
        <p id="welcome-message" className="text-lg">
          Welcome to the Austin Pinball Collective!
        </p>
        <p id="sign-in-message" className="text-md italic">
          Please sign in below
        </p>
      </div>
      <div
        id="sign-in-form"
        className="border-l-4 border-b-2 p-2 rounded-sm bg-slate-100 flex flex-col gap-2"
      >
        <div id="personal-info" className="flex flex-col gap-2">
          <label id="full-name-label">Full Name</label>
          <input
            id="full-name-input"
            type="text"
            className="border-2 rounded-sm md:w-1/2 px-2 py-1"
            placeholder="Keith Elwin"
            onInput={(input) => {
              const newName = getInputText(input);
              setName(newName);
              updateAgreementState(
                setAgreementState,
                checkAgreeDebouncing,
                newName,
                email,
              );
            }}
          />
          <label id="email-label">Email</label>
          <input
            id="email-input"
            type="email"
            className="border-2 rounded-sm md:w-1/2 px-2 py-1"
            placeholder="keith.elwin@hotmail.com"
            onInput={(input) => {
              const newEmail = getInputText(input);
              setEmail(getInputText(input));
              updateAgreementState(
                setAgreementState,
                checkAgreeDebouncing,
                name,
                newEmail,
              );
            }}
          />
        </div>
        {agreementState === AlreadyAgreed
          ? <div>Welcome Back!</div>
          : (
            <AgreementBlock
              onOpenAgreement={() => setAgreementOpen(true)}
              onCheckboxChange={(checked: boolean) =>
                checked
                  ? setAgreementState(AgreementCompleted)
                  : setAgreementState(AgreementNeeded)}
            />
          )}
        <div>
          {submitEnabled
            ? (
              <button
                type="submit"
                className="bg-emerald-400 not-active:border-l-2 border-emerald-700 rounded-sm px-4 py-1 text-lg hover:cursor-pointer hover:bg-emerald-500"
                onClick={() => {
                  window.location.href = "/signin-success";
                }}
              >
                Submit
              </button>
            )
            : (
              <button
                type="submit"
                className="bg-slate-200 border-slate-300 border-l-2 rounded-sm px-4 py-1 text-lg text-slate-50 hover:cursor-not-allowed"
                disabled
              >
                Submit
              </button>
            )}
        </div>
      </div>
      {agreementOpen
        ? (
          <Modal>
            <div
              id="modal-page"
              className="m-2 md:m-10 bg-slate-200 opacity-100 rounded-sm"
            >
              <div
                id="agreement-modal-header"
                className="flex flex-row justify-between p-2"
              >
                <p className="text-xl">
                  Austin Pinball Collective Member Agreement
                </p>
                <button
                  id="agreement-exit"
                  type="button"
                  className="bg-slate-300 rounded-lg px-2 py-1"
                  onClick={() => setAgreementOpen(false)}
                >
                  Close
                </button>
              </div>
              <div id="agreement-content" className="p-2">
                <p className="text-md">Blah blah blah blah</p>
              </div>
            </div>
          </Modal>
        )
        : <div></div>}
      <Link
        id="signin-success-link"
        className="invisible"
        href="/signin-success"
      />
    </div>
  );
}
