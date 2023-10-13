import { selectIsLoggedIn } from "@/feature/store/auth/auth.selectors";
import { List, ListItemButton, Typography } from "@mui/joy";
import { useAppSelector } from "@/feature/store/hooks";
import { useNavigate } from "react-router-dom";

export enum NavItem {
  New = "New",
  Popular = "Popular",
  Subscriptions = "Subscriptions",
}

export interface NavMenuProps {
  selectedItem?: NavItem;
}

export default function NavMenu({ selectedItem }: NavMenuProps) {
  const nav = useNavigate();

  const isLoggedIn = useAppSelector(selectIsLoggedIn);

  const items = [
    {
      label: NavItem.New,
      onClick: () =>
        nav(`/?videos=${NavItem.New.toLowerCase()}`, { replace: true }),
    },
    {
      label: NavItem.Popular,
      onClick: () =>
        nav(`/?videos=${NavItem.Popular.toLowerCase()}`, { replace: true }),
    },
  ];

  if (isLoggedIn)
    items.push({
      label: NavItem.Subscriptions,
      onClick: () =>
        nav(`/?videos=${NavItem.Subscriptions.toLowerCase()}`, {
          replace: true,
        }),
    });

  return (
    <List style={{ gap: 5 }}>
      {items.map((item) => (
        <ListItemButton
          key={item.label}
          selected={selectedItem !== undefined && item.label === selectedItem}
          onClick={item.onClick}
        >
          <Typography level="title-lg" fontSize="1em">
            {item.label}
          </Typography>
        </ListItemButton>
      ))}
    </List>
  );
}
