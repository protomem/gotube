export default function Page({ params }: { params: { nickname: string } }) {
  return <div>Profile {params.nickname}</div>;
}
