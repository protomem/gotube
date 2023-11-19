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
import { formatViews } from "@/lib";
import Link from "next/link";
import { ROUTES } from "@/shared/constants/routes";

interface VideoCardProps {
  video: VideoEntity;
}

const cardContentImageLoader: ImageLoader = ({ src, width, quality }) => {
  return `https://images.unsplash.com/${src}?w=${width}&dpr=2&q=${quality}`;
};

export function VideoCard({ video }: VideoCardProps) {
  return (
    <Card className="min-w-[290px] w-1/3">
      <CardContent className="p-0">
        <AspectRatio ratio={16 / 9} className="bg-muted">
          <Image
            loader={cardContentImageLoader}
            src="photo-1588345921523-c2dcdb7f1dcd"
            alt="Photo by Drew Beamer"
            fill
            className="rounded-md object-cover"
          />
        </AspectRatio>
      </CardContent>

      <CardHeader>
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
      </CardHeader>
    </Card>
  );
}
