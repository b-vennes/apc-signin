"use client";

import Link from "next/link";
import Modal from "@/components/modal.tsx";
import { useState } from "react";

export default function SignIn() {
  const [agreementOpen, setAgreementOpen] = useState(false);

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
            className="border-2 rounded-sm md:w-1/2 bg-slate-200 focus:bg-slate-50"
          />
          <label id="email-label">Email</label>
          <input
            id="email-input"
            type="email"
            className="border-2 rounded-sm md:w-1/2 bg-slate-200 focus:bg-slate-50"
          />
        </div>
        <div
          id="member-agreement-info"
          className="flex flex-col items-start"
        >
          <div
            id="member-agreement-info"
            className="flex flex-row justify-between gap-4"
          >
            <p id="member-agreement-title">Member Agreement</p>
            <button
              type="button"
              id="view-agreement-button"
              className="bg-slate-500 rounded-sm px-4 py-1"
              onClick={() => {
                setAgreementOpen(true);
                console.log("open agreement clicked");
              }}
            >
              View
            </button>
          </div>
          <div className="flex flex-row gap-1">
            <input id="member-agreement-checkbox" type="checkbox" />
            <label id="agree-checkbox" className="text-sm">
              I agree to the rules of the Austin Pinball Collective space set
              forth in the Austin Pinball Collective Member Agreement
            </label>
          </div>
        </div>
        <div>
          <button
            type="submit"
            className="bg-slate-500 rounded-sm px-4 py-1 text-lg"
            onClick={() => {
              window.location.href = "/signin-success";
            }}
          >
            Submit
          </button>
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
