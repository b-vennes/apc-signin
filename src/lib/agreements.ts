import { useState } from "react";

export const AgreementUndecided = "agreement-undecided";

export const AgreementNeeded = "agreement-needed";

export const AgreementCompleted = "agreement-completed";

export const AlreadyAgreed = "already-agreed";

export type AgreementState =
  | typeof AgreementUndecided
  | typeof AgreementNeeded
  | typeof AgreementCompleted
  | typeof AlreadyAgreed;

export function defaultAgreementState(): AgreementState {
  return AgreementUndecided;
}

export function useAgreement(): [
  AgreementState,
  (state: AgreementState) => void,
] {
  return useState(defaultAgreementState);
}
