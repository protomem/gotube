import numeral from "numeral";

export function formatViews(views: number): string {
  if (views < 1000) return numeral(views).format("0 a");
  return numeral(views).format("0.0 a");
}
