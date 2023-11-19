import { videos } from "@/fixtures/videos";
import { VideoEntity } from "@/entities/domain/models";

interface GetVideoByIdRequest {
  id: string;
}

interface GetVideoByIdResponse {
  video: VideoEntity;
}

class VideoService {
  constructor() {}

  getById(_: GetVideoByIdRequest) {
    return { video: videos[0] } as GetVideoByIdResponse;
  }
}

export const videoService = new VideoService();
