export enum ROUTES {
  HOME = "/",
  SEARCH = "/search",
}

export function withQuery(basePath: string, query: Record<string, string>) {
  return `${basePath}?${new URLSearchParams(query).toString()}`;
}
