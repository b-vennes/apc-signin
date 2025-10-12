import Link from "next/link";

export default function SignInSuccess() {
  return (
    <div id="sign-in-success-root">
      <p>Success! Have fun!</p>
      <Link href="/signin">Back To Sign In</Link>
    </div>
  );
}
