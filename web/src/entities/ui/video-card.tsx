import { VideoEntity } from "@/entities/domain/models";

import Image, { ImageLoader } from "next/image";
import { AspectRatio } from "@/shared/ui/aspect-ratio";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/card";
import { cn, formatViews } from "@/lib";
import Link from "next/link";
import { ROUTES } from "@/shared/constants/routes";
import dynamic from "next/dynamic";

const DynamicAvatar = dynamic(() => import("@/shared/components/avatar"), {
  ssr: false,
});

interface VideoCardProps {
  video: VideoEntity;
  composit?: "horizontal" | "vertical";
}

const cardContentImageLoader: ImageLoader = ({ src, width, quality }) => {
  return `https://images.unsplash.com/${src}?w=${width}&dpr=2&q=${quality}`;
};

export function VideoCard({ video, composit }: VideoCardProps) {
  if (composit === undefined) composit = "vertical";

  const verticalCardStyles = cn("min-w-[290px] w-1/3 flex flex-col");
  const horintalCardStyles = cn(
    "min-w-[700px] w-auto flex flex-row justify-start items-center",
  );

  const verticalCardContentStyles = cn("p-0");
  const horintalCardContentStyles = cn("p-1 w-[300px]");

  const verticalCardHeaderStyles = cn("flex flex-row gap-3 p-3");
  const horintalCardHeaderStyles = cn(
    "flex flex-col justify-between gap-3 p-0 pl-5 pt-3 self-start",
  );

  return (
    <Card
      className={
        composit === "vertical" ? verticalCardStyles : horintalCardStyles
      }
    >
      <CardContent
        className={
          composit === "vertical"
            ? verticalCardContentStyles
            : horintalCardContentStyles
        }
      >
        <AspectRatio ratio={16 / 10} className="bg-muted">
          <Image
            loader={cardContentImageLoader}
            src="photo-1588345921523-c2dcdb7f1dcd"
            alt="Photo by Drew Beamer"
            fill
            className="rounded-md object-cover"
          />
        </AspectRatio>
      </CardContent>

      <CardHeader
        className={
          composit === "vertical"
            ? verticalCardHeaderStyles
            : horintalCardHeaderStyles
        }
      >
        {composit === "vertical" ? (
          <>
            <div className="pt-2">
              <DynamicAvatar title={video.author.nickname.slice(0, 2)} />
            </div>
            <div>
              <CardTitle className="text-xl">
                <Link href={`${ROUTES.WATCH}/${video.id}`}>{video.title}</Link>
              </CardTitle>
              <CardDescription className="text-md">
                <Link href={`${ROUTES.PROFILE}/${video.author.nickname}`}>
                  {video.author.nickname}
                </Link>
              </CardDescription>
              <CardDescription>
                {formatViews(video.views + 1000)} views •
                {" " + video.updatedAt.toDateString()}
              </CardDescription>
            </div>
          </>
        ) : (
          <>
            <div className="flex flex-col gap-[0.2rem]">
              <CardTitle className="text-xl">
                <Link href={`${ROUTES.WATCH}/${video.id}`}>{video.title}</Link>
              </CardTitle>
              <CardDescription>
                {formatViews(video.views + 1000)} views •
                {" " + video.updatedAt.toDateString()}
              </CardDescription>
              <CardDescription className="text-md flex flex-row justify-start items-center gap-2">
                <DynamicAvatar
                  title={video.author.nickname.slice(0, 2)}
                  width={35}
                  height={35}
                />
                <Link href={`${ROUTES.PROFILE}/${video.author.nickname}`}>
                  {video.author.nickname}
                </Link>
              </CardDescription>
            </div>
            <div>
              <CardDescription>{video.description}</CardDescription>
            </div>
          </>
        )}
      </CardHeader>
    </Card>
  );
}
