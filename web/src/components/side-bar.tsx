import { useNavigate } from "react-router-dom";
import NavMenu, { NavMenuItem } from "./nav-menu";
import { Box } from "@chakra-ui/react";

type Props = {
  type?: "minimal" | "expanded";
  selectedNavMenuItem?: NavMenuItem;
};

const SideBar = ({ type, selectedNavMenuItem }: Props) => {
  type = type || "expanded";

  const nav = useNavigate();
  const handleSelect = (item: NavMenuItem) => {
    nav(`/?nav=${item.toLowerCase()}`, { replace: true });
  };

  return (
    <Box width={type === "minimal" ? "4rem" : "10rem"} paddingLeft="4">
      <NavMenu
        type={type}
        selectedItem={selectedNavMenuItem}
        items={[
          NavMenuItem.Home,
          NavMenuItem.Trends,
          NavMenuItem.Subscriptions,
        ]}
        onItemSelect={handleSelect}
      />

      {/* <Divider my="2" /> */}
    </Box>
  );
};

export default SideBar;
