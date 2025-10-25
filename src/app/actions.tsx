"use server";

export async function memberAgreed(
  name: string,
  email: string,
): Promise<boolean> {
  const response = await fetch(
    "http://localhost:3339/api/member/agreed",
    {
      method: "POST",
      body: `{"name": "${name}", "email": "${email}"}`,
    },
  );

  const body = await response.text();

  return body === "true";
}
