import { useNavigate } from "react-router-dom";
import { useAuth } from "../providers/auth-provider";
import NavMenu, { NavMenuItem } from "./nav-menu";
import { Box, Divider } from "@chakra-ui/react";
import SubscriptionList from "./subscriptions-list";

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

  const { isAuthenticated, currentUser } = useAuth();

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

      {isAuthenticated && currentUser && (
        <>
          <Divider my="4" />

          <SubscriptionList type={type} user={currentUser} />
        </>
      )}
    </Box>
  );
};

export default SideBar;
