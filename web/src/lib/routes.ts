export enum ROUTES {
  HOME = "/",
  SEARCH = "/search",
  WATCH = "/watch",
  PROFILE = "/profile",
}

export function withQuery(basePath: string, query: Record<string, string>) {
  return `${basePath}?${new URLSearchParams(query).toString()}`;
}
