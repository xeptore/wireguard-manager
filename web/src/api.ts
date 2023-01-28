import { API_BASE_URL } from "@/constants";
import { URL } from "node:url";

export function path(route: string): string {
  const url = new URL(route, API_BASE_URL);

  return url.href;
}
