import numeral from "numeral";

export function capitalize(str: string) {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

export function repeat<T>(array: T[], size: number) {
  let result: T[] = [];
  for (let i = 0; i < size; i++) result.push(...array);
  return result;
}

export function formatDate(date: Date) {
  const now = new Date();
  const diff = now.getTime() - date.getTime();

  if (diff < 1000 * 60) return "now";
  if (diff < 1000 * 60 * 60) return `${Math.floor(diff / (1000 * 60))}m`;
  if (diff < 1000 * 60 * 60 * 24)
    return `${Math.floor(diff / (1000 * 60 * 60))} hours`;
  if (diff < 1000 * 60 * 60 * 24 * 7) {
    if (date.getDate() === now.getDate() - 1) return "yesterday";

    return `${Math.floor(diff / (1000 * 60 * 60 * 24))} days`;
  }
  if (diff < 1000 * 60 * 60 * 24 * 30)
    return `${Math.floor(diff / (1000 * 60 * 60 * 24 * 7))} weak`;
  if (diff < 1000 * 60 * 60 * 24 * 365)
    return `${Math.floor(diff / (1000 * 60 * 60 * 24 * 30))} month`;
  return `${Math.floor(diff / (1000 * 60 * 60 * 24 * 365))} years`;
}

export function formatViews(views: number) {
  if (views < 1000) return views;
  if (!(views % 1000)) return numeral(views).format("0a");
  return numeral(views).format("0.0a");
}
