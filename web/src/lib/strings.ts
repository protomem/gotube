export function capitalize(string: string): string {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

export function appendQuery(
  base: string,
  ...queries: { key: string; value: string }[]
): string {
  return `${base}?${queries.map((q) => `${q.key}=${q.value}`).join("&")}`;
}
