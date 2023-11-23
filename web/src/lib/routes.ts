export enum ROUTES {
  HOME = "/",
}

export function withQuery(basePath: string, query: Record<string, string>) {
  return `${basePath}?${new URLSearchParams(query).toString()}`;
}
