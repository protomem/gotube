import { selectIsLoggedIn } from "@/feature/store/auth/auth.selectors";
import { useAppDispatch, useAppSelector } from "@/feature/store/hooks";
import { selectSelectedItem } from "@/feature/store/nav/nav.selectors";
import { NavItem, navActions } from "@/feature/store/nav/nav.slice";
import { List, ListItemButton, Typography } from "@mui/joy";

export default function NavMenu() {
  const dispatch = useAppDispatch();

  const isLoggedIn = useAppSelector(selectIsLoggedIn);
  const selectedItem = useAppSelector(selectSelectedItem);

  const items = [
    {
      label: NavItem.New,
      onClick: () => dispatch(navActions.setSelectedItem(NavItem.New)),
    },
    {
      label: NavItem.Popular,
      onClick: () => dispatch(navActions.setSelectedItem(NavItem.Popular)),
    },
  ];

  if (isLoggedIn)
    items.push({
      label: NavItem.Subscriptions,
      onClick: () =>
        dispatch(navActions.setSelectedItem(NavItem.Subscriptions)),
    });

  return (
    <List style={{ gap: 5 }}>
      {items.map((item) => (
        <ListItemButton
          key={item.label}
          onClick={item.onClick}
          selected={item.label === selectedItem}
        >
          <Typography level="title-lg" fontSize="1em">
            {item.label}
          </Typography>
        </ListItemButton>
      ))}
    </List>
  );
}
