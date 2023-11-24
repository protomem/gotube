export function capitalize(str: string) {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

export function repeat<T>(array: T[], size: number) {
  let result: T[] = [];
  for (let i = 0; i < size; i++) result.push(...array);
  return result;
}
