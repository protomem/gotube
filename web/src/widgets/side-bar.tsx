import { NavMenu, NavMenuProps } from "@/feature/nav-menu";

interface SideBarProps extends NavMenuProps {}

export function SideBar({ ...props }: SideBarProps) {
  return (
    <div className="px-3 py-3">
      <NavMenu {...props} />
    </div>
  );
}

export { Navigates } from "@/feature/nav-menu";
