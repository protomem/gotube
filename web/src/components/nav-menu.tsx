import _ from "lodash";
import { JSX } from "react";
import NextLink from "next/link";
import { FaFire, FaHouse } from "react-icons/fa6";
import { Button, List, ListItem } from "@chakra-ui/react";

export type NavItem = {
  icon: JSX.Element;
  label: string;
  href: string;
  selected: boolean;
};

export const defaultNavItems: NavItem[] = [
  {
    icon: <FaHouse />,
    label: "Home",
    href: "/?videos=home",
    selected: false,
  },
  {
    icon: <FaFire />,
    label: "Trends",
    href: "/?videos=trends",
    selected: false,
  },
];

interface Props {
  navItems?: NavItem[];
  labelSelected?: string;
}

export default function NavMenu({
  labelSelected,
  navItems = defaultNavItems,
}: Props) {
  navItems = _.uniqBy(navItems, (item) => item.label);
  const chooseSelectedLabel = (
    currentLabel: string,
    currentLabelSelected: boolean,
    labelSelected?: string,
  ): boolean => {
    if (
      labelSelected !== undefined &&
      labelSelected.toLowerCase() === currentLabel.toLowerCase()
    )
      return true;
    return currentLabelSelected;
  };

  return (
    <List spacing="2" w="100%">
      {navItems.map((item) => (
        <ListItem key={item.label}>
          <Button
            w="100%"
            as={NextLink}
            href={item.href}
            leftIcon={item.icon}
            variant={
              chooseSelectedLabel(item.label, item.selected, labelSelected)
                ? "solid"
                : "ghost"
            }
            justifyContent="start"
            gap="2"
          >
            {item.label}
          </Button>
        </ListItem>
      ))}
    </List>
  );
}
