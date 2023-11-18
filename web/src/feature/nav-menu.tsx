import React from "react";
import { capitalize } from "@/lib";
import { ROUTES } from "@/shared/constants/routes";

import { Button } from "@/shared/ui/button";
import Link from "next/link";

export enum Navigates {
  New = "new",
  Popular = "popular",
  Subscriptions = "subscriptions",
}

interface NavMenuItemProps {
  title: string;
  selected?: boolean;
  href: string;
}

const NavMenuItem: React.FC<NavMenuItemProps> = ({ title, selected, href }) => {
  if (selected === undefined) selected = false;

  return (
    <Button
      asChild
      disabled={!selected}
      variant={selected ? "default" : "ghost"}
      className="justify-start"
    >
      <Link href={href}>{capitalize(title)}</Link>
    </Button>
  );
};

export interface NavMenuProps {
  selectedNav?: Navigates;
  withSubscriptions?: boolean;
}

export function NavMenu({ selectedNav, withSubscriptions }: NavMenuProps) {
  if (selectedNav === undefined) selectedNav = Navigates.New;
  if (withSubscriptions === undefined) withSubscriptions = false;

  const navMenuItems = [
    {
      title: Navigates.New,
      selected: selectedNav === Navigates.New,
      href: ROUTES.HOME,
    },
    {
      title: Navigates.Popular,
      selected: selectedNav === Navigates.Popular,
      href: ROUTES.HOME,
    },
  ];

  if (withSubscriptions)
    navMenuItems.push({
      title: Navigates.Subscriptions,
      selected: selectedNav === Navigates.Subscriptions,
      href: ROUTES.HOME,
    });

  return (
    <div className="flex flex-col gap-3">
      {navMenuItems.map((item, i) => (
        <NavMenuItem
          key={i}
          title={item.title}
          selected={item.selected}
          href={item.href}
        />
      ))}
    </div>
  );
}
