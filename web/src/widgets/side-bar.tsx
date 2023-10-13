import { Box, Divider } from "@mui/joy";
import NavMenu, { NavItem } from "@/feature/nav-menu";
import { selectIsLoggedIn } from "@/feature/store/auth/auth.selectors";
import { useAppSelector } from "@/feature/store/hooks";
import SubscriptionsList from "@/feature/subscriptions-list";

export interface SideBarProps {
  selectedNavItem?: NavItem;
}

export default function SideBar({ selectedNavItem }: SideBarProps) {
  const isLoggedIn = useAppSelector(selectIsLoggedIn);

  return (
    <Box
      style={{
        marginTop: 20,
        marginLeft: 10,
        marginRight: 10,
      }}
    >
      <NavMenu selectedItem={selectedNavItem} />

      <Divider />

      {isLoggedIn && <SubscriptionsList />}
    </Box>
  );
}
