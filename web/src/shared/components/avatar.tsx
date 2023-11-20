import Image, { ImageLoader } from "next/image";

function generateAvatar(text: string, fgcolor: string, bgcolor: string) {
  const canvas = document.createElement("canvas");
  const ctx = canvas.getContext("2d");
  if (ctx == null) throw new Error("Canvas is not supported");

  canvas.width = 200;
  canvas.height = 200;

  ctx.fillStyle = bgcolor;
  ctx.fillRect(0, 0, canvas.width, canvas.height);

  ctx.font = "bold 70px Assistant";
  ctx.fillStyle = fgcolor;
  ctx.textAlign = "center";
  ctx.textBaseline = "middle";
  ctx.fillText(text, canvas.width / 2, canvas.height / 2);

  return canvas.toDataURL("image/png");
}

function colorFromText(text: string) {
  let hash = 0;
  for (let i = 0; i < text.length; i++) {
    hash = text.charCodeAt(i) + ((hash << 5) - hash);
  }
  let color = "#";
  for (let i = 0; i < 3; i++) {
    const value = (hash >> (i * 8)) & 0xff;
    color += ("00" + value.toString(16)).substr(-2);
  }
  return color;
}

interface AvatarProps {
  src?: string;
  title?: string;
  loader?: ImageLoader;
  alt?: string;
}

export function Avatar({ title, alt }: AvatarProps) {
  return (
    <Image
      src={generateAvatar(title ?? "", "white", colorFromText(title ?? ""))}
      alt={alt ?? "avatar"}
      width={50}
      height={50}
      className="rounded-full"
    />
  );
}

export default Avatar;
