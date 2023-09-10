export default function Page({ params }: { params: { videoId: string } }) {
  return <div>Video: {params.videoId}</div>;
}
