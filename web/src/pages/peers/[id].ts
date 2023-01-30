import type { APIRoute } from "astro";
import httpStatusCodes from "http-status-codes";
import pWaterfall from "p-waterfall";
import QRCode from "qrcode-svg";

import { path } from "@/api";

export const get: APIRoute = async function get({ params, request, redirect }) {
  const { id } = params;

  return await pWaterfall([
    async () => {
      try {
        return await fetch(path(`peers/${id}`), {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Cookie: request.headers.get("cookie") ?? "",
          },
          redirect: "follow",
          referrerPolicy: "no-referrer",
          credentials: "include",
          body: null,
        });
      } catch (error) {
        console.error(error);
        return await Promise.resolve(new Response());
      }
    },
    async (apiRes) => {
      if (apiRes.status === httpStatusCodes.OK) {
        const content = await apiRes.text();
        const svg = new QRCode(content).svg();
        return {
          body: Buffer.from(svg).toString("base64"),
          headers: {
            ...Object.fromEntries(apiRes.headers.entries()),
            "Content-Type": "image/svg+xml",
          },
        };
      }

      if (apiRes.status === httpStatusCodes.UNAUTHORIZED) {
        return redirect(
          "/?error=401",
          httpStatusCodes.TEMPORARY_REDIRECT as 307
        );
      }

      return redirect("/?error=500", httpStatusCodes.TEMPORARY_REDIRECT as 307);
    },
  ]);
};
