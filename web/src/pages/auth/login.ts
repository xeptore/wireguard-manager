import type { APIRoute } from "astro";
import httpStatusCodes from "http-status-codes";
import { path } from "@/api";

export const post: APIRoute = async function post({ request, redirect }) {
  const formData = await request.formData();
  const values = Object.fromEntries(formData.entries());

  const apiResponse = await fetch(path("auth/login"), {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    redirect: "follow",
    referrerPolicy: "no-referrer",
    credentials: "same-origin",
    body: JSON.stringify(values),
  });

  if (apiResponse.status === httpStatusCodes.CREATED) {
    const response = redirect("/", httpStatusCodes.TEMPORARY_REDIRECT as 301);
    Array.from(apiResponse.headers.entries()).forEach(([k, v]) => response.headers.set(k, v));

    return response;
  }

  if (apiResponse.status === httpStatusCodes.UNAUTHORIZED) {
    return redirect("/login?error=401", httpStatusCodes.TEMPORARY_REDIRECT as 301);
  }

  return redirect("/login?error=500", httpStatusCodes.TEMPORARY_REDIRECT as 301);
};
