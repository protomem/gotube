import { Video } from "@/entities/models";
import {
  AspectRatio,
  Button,
  Card,
  CardContent,
  CardOverflow,
  Divider,
  Link,
  Typography,
} from "@mui/joy";
import { Link as RouterLink, useNavigate } from "react-router-dom";
import Avatar from "@/shared/avatar";

export interface VideoCardProps {
  video: Video;
}

export default function VideoCard({ video }: VideoCardProps) {
  const nav = useNavigate();

  return (
    <Card variant={"outlined"} style={{ minWidth: "220px", maxWidth: "340px" }}>
      <CardOverflow>
        <AspectRatio ratio={2}>
          <Button component={RouterLink} to={`/video/${video.id}`}>
            <img
              src={video.thumbnailPath}
              alt={video.title}
              style={{ objectFit: "cover" }}
            />
          </Button>
        </AspectRatio>
      </CardOverflow>

      <CardContent>
        <Link
          component={RouterLink}
          to={`/video/${video.id}`}
          level={"title-md"}
          underline={"none"}
        >
          {video.title}
        </Link>

        <Typography level={"body-md"} noWrap>
          {video.description}
        </Typography>
      </CardContent>

      <CardOverflow
        variant={"plain"}
        style={{ backgroundColor: "background.level1" }}
      >
        <Divider inset="context" />

        <CardContent orientation="horizontal" style={{ alignItems: "center" }}>
          <Button
            variant={"plain"}
            color={"neutral"}
            style={{
              display: "flex",
              flexDirection: "row",
              alignItems: "center",
              gap: 6,
            }}
            onClick={() => {
              nav(`/profile/${video.user.nickname}`);
            }}
          >
            <Avatar user={video.user} size={"sm"} />
            <Typography
              width={"30px"}
              level={"body-xs"}
              fontWeight={"md"}
              textColor={"text.secondary"}
              noWrap
            >
              {video.user.nickname}
            </Typography>
          </Button>

          <Divider orientation={"vertical"} />

          <Typography
            level={"body-xs"}
            fontWeight={"md"}
            textColor={"text.secondary"}
          >
            {video.views} views
          </Typography>
        </CardContent>
      </CardOverflow>
    </Card>
  );
}
