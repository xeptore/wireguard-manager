import type { APIRoute } from "astro";
import httpStatusCodes from "http-status-codes";
import pWaterfall from "p-waterfall";
import { path } from "@/api";

export const post: APIRoute = async function post({ request, redirect }) {
  return await pWaterfall(
    [
      async (req) => await req.formData(),
      async (formData) => {
        const values = Object.fromEntries(formData.entries());
        try {
          return await fetch(path("auth/login"), {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            redirect: "follow",
            referrerPolicy: "no-referrer",
            credentials: "same-origin",
            body: JSON.stringify(values),
          });
        } catch (error) {
          console.error(error);
          return await Promise.resolve(new Response());
        }
      },
      (apiRes) => {
        if (apiRes.status === httpStatusCodes.CREATED) {
          const response = redirect(
            "/",
            httpStatusCodes.TEMPORARY_REDIRECT as 307
          );
          Array.from(apiRes.headers.entries()).forEach(([k, v]) =>
            response.headers.set(k, v)
          );

          return response;
        }

        if (apiRes.status === httpStatusCodes.UNAUTHORIZED) {
          return redirect(
            "/login?error=401",
            httpStatusCodes.TEMPORARY_REDIRECT as 307
          );
        }

        return redirect(
          "/login?error=500",
          httpStatusCodes.TEMPORARY_REDIRECT as 307
        );
      },
    ],
    request
  );
};
