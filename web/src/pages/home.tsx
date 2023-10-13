import { NavItem } from "@/feature/nav-menu";
import AppBar from "@/widgets/app-bar";
import Layout from "@/widgets/layout";
import SideBar from "@/widgets/side-bar";
import { Typography } from "@mui/joy";
import { useSearchParams } from "react-router-dom";

export default function Home() {
  const [searchParams] = useSearchParams();

  let selectedNavItem: NavItem;
  switch (searchParams.get("videos")) {
    case NavItem.New.toLowerCase():
      selectedNavItem = NavItem.New;
      break;
    case NavItem.Popular.toLowerCase():
      selectedNavItem = NavItem.Popular;
      break;
    case NavItem.Subscriptions.toLowerCase():
      selectedNavItem = NavItem.Subscriptions;
      break;
    default:
      selectedNavItem = NavItem.New;
      break;
  }

  return (
    <Layout>
      <AppBar />

      <SideBar selectedNavItem={selectedNavItem} />

      <Typography
        fontSize={25}
        fontWeight={"bold"}
        textAlign={"center"}
        style={{ marginTop: 20 }}
      >
        {selectedNavItem}
      </Typography>
    </Layout>
  );
}
