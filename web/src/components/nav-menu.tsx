import _ from "lodash";
import { JSX, useRef } from "react";
import { useSize } from "@chakra-ui/react-use-size";
import NextLink from "next/link";
import { FaFire, FaHouse } from "react-icons/fa6";
import { Button, IconButton, List, ListItem } from "@chakra-ui/react";

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

  const ref = useRef<HTMLUListElement>(null);
  const size = useSize(ref);

  const isShort = size && size.width && size.width < 100;

  return (
    <List ref={ref} spacing="2" w="100%">
      {navItems.map((item) => (
        <ListItem key={item.label}>
          {isShort ? (
            <IconButton
              aria-label={`Icon ${item.label}`}
              as={NextLink}
              href={item.href}
              icon={item.icon}
              variant={
                chooseSelectedLabel(item.label, item.selected, labelSelected)
                  ? "solid"
                  : "ghost"
              }
            />
          ) : (
            <Button
              w="100%"
              p="3"
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
          )}
        </ListItem>
      ))}
    </List>
  );
}
