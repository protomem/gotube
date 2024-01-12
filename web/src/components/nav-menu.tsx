import { FaHome, FaMeteor, FaUserFriends } from "react-icons/fa";
import { Button, Text, VStack } from "@chakra-ui/react";

export enum NavMenuItem {
  Home = "Home",
  Trends = "Trends",
  Subscriptions = "Subscriptions",
}

export const resolveNavMenuItem = (s: string | null) => {
  switch (s) {
    case "popular":
    case NavMenuItem.Trends.toLowerCase():
      return NavMenuItem.Trends;

    case "subs":
    case NavMenuItem.Subscriptions.toLowerCase():
      return NavMenuItem.Subscriptions;

    default:
      return NavMenuItem.Home;
  }
};

const NavMenuItemIcons = {
  [NavMenuItem.Home]: <FaHome />,
  [NavMenuItem.Trends]: <FaMeteor />,
  [NavMenuItem.Subscriptions]: <FaUserFriends />,
};

type Props = {
  type?: "minimal" | "expanded";
  selectedItem?: NavMenuItem;
  items?: NavMenuItem[];
  onItemSelect?: (item: NavMenuItem) => void;
};

const NavMenu = ({ type, selectedItem, items, onItemSelect }: Props) => {
  return (
    <VStack align="start">
      {items?.map((item) => (
        <Button
          key={item}
          w="full"
          size="md"
          variant={selectedItem === item ? "solid" : "ghost"}
          colorScheme={selectedItem === item ? "teal" : "gray"}
          leftIcon={type === "expanded" ? NavMenuItemIcons[item] : undefined}
          onClick={() => onItemSelect?.(item)}
          alignContent="start"
          sx={{ justifyContent: "start" }}
        >
          {type === "expanded" ? (
            <Text fontSize="sm">{item}</Text>
          ) : (
            NavMenuItemIcons[item]
          )}
        </Button>
      ))}
    </VStack>
  );
};

export default NavMenu;
