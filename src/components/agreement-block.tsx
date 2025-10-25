"use client";

function getInputCheckbox(
  event: unknown,
): boolean {
  const eventWithChecked = event as { target: { checked: boolean } };

  return eventWithChecked?.target?.checked ?? false;
}

/**
  Properties for the AgreementBlock component.

  @property {() => void} onOpenAgreeement a function which executes when a
    person wants to open the agreement

  @property {(boolean) => void} onCheckboxChange a function which executes when
    a person has changed their member agreement decision
 */
export type AgreementBlockProps = {
  onOpenAgreement: () => void;
  onCheckboxChange: (checked: boolean) => void;
};

/**
  Defines a block which lets people trigger a command to open and approve the
  member agreement.
 */
export function AgreementBlock(props: AgreementBlockProps) {
  return (
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
          className="not-active:border-l-1 bg-amber-100 border-amber-400 rounded-sm px-4 py-1 hover:cursor-pointer hover:bg-amber-200"
          onClick={props.onOpenAgreement}
        >
          View
        </button>
      </div>
      <div className="flex flex-row gap-1">
        <input
          id="member-agreement-checkbox"
          type="checkbox"
          onClick={(event) => props.onCheckboxChange(getInputCheckbox(event))}
        />
        <label
          id="agree-checkbox"
          className="text-sm"
        >
          I agree to the rules of the Austin Pinball Collective space set forth
          in the Austin Pinball Collective Member Agreement
        </label>
      </div>
    </div>
  );
}
