export function repeat<T>(array: T[], n?: number): T[] {
  if (n === undefined) n = 1;
  let repeated: T[] = [];
  for (let i = 0; i < n; i++) {
    repeated = [...repeated, ...array];
  }
  return repeated;
}

export function chunk<T>(array: T[], n?: number): T[][] {
  if (n === undefined) n = 3;
  let chunks: T[][] = [];
  for (let i = 0; i < array.length; i += n) {
    chunks.push(array.slice(i, i + n));
  }
  return chunks;
}
