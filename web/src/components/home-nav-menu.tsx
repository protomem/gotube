import React from "react";
import { ROUTES, withQuery } from "@/lib/routes";

import NextLink from "next/link";
import { Button, VStack } from "@chakra-ui/react";
import { capitalize } from "@/lib/utils";

export enum HomeNavMenuItemLabel {
  New = "new",
  Popular = "popular",
  Subscriptions = "subscriptions",
}

type HomeNavMenuItemProps = {
  label: HomeNavMenuItemLabel;
  href: string;
  selected: boolean;
};

const HomeNavMenuItem = ({ label, href, selected }: HomeNavMenuItemProps) => {
  return (
    <Button
      as={NextLink}
      href={href}
      colorScheme={selected ? "teal" : "gray"}
      variant={selected ? "solid" : "ghost"}
      sx={{ flexDirection: "column", alignItems: "start" }}
    >
      {capitalize(label)}
    </Button>
  );
};

type HomeNavMenuProps = {
  selectedItem?: HomeNavMenuItemLabel;
  withSubscriptions?: boolean;
};

const HomeNavMenu = ({ selectedItem, withSubscriptions }: HomeNavMenuProps) => {
  const itemsProps: HomeNavMenuItemProps[] = [
    {
      label: HomeNavMenuItemLabel.New,
      href: withQuery(ROUTES.HOME, { nav: "new" }),
      selected: selectedItem === HomeNavMenuItemLabel.New,
    },
    {
      label: HomeNavMenuItemLabel.Popular,
      href: withQuery(ROUTES.HOME, { nav: "popular" }),
      selected: selectedItem === HomeNavMenuItemLabel.Popular,
    },
  ];

  if (withSubscriptions)
    itemsProps.push({
      label: HomeNavMenuItemLabel.Subscriptions,
      href: withQuery(ROUTES.HOME, { nav: "subscriptions" }),
      selected: selectedItem === HomeNavMenuItemLabel.Subscriptions,
    });

  return (
    <VStack align="stretch">
      {itemsProps.map((props, i) => (
        <HomeNavMenuItem key={i} {...props} />
      ))}
    </VStack>
  );
};

export default HomeNavMenu;
