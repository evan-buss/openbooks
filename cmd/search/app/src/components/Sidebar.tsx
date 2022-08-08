import { Navbar, NavLink, useMantineTheme } from "@mantine/core";
import { ChartLineUp, Gear, SignOut, UserCircle } from "phosphor-react";
import { Link, useLocation } from "react-router-dom";
import useAuth from "../util/auth";

interface SidebarProps {
  opened: boolean;
  setOpened: (opened: boolean) => void;
}

export function Sidebar({ opened }: SidebarProps) {
  const { logOut } = useAuth();
  const { user } = useAuth();

  return (
    <Navbar p="md" hiddenBreakpoint="sm" hidden={!opened} width={{ sm: 250 }}>
      <Navbar.Section grow mt="xs">
        <NavBarLink
          path="/"
          label="Dashboard"
          icon={<ChartLineUp size={16} />}
        />
        <NavBarLink
          path="/settings"
          label="Settings"
          icon={<Gear size={16} />}
        />
      </Navbar.Section>
      <Navbar.Section>
        <NavBarLink label={user?.name} icon={<UserCircle size={16} />} />
        <NavBarLink
          label="Log Out"
          icon={<SignOut size={16} />}
          onClick={logOut}
        />
      </Navbar.Section>
    </Navbar>
  );
}

interface NavBarLinkProps {
  label?: string;
  icon: JSX.Element;
  path?: string;
  onClick?: () => void;
}

function NavBarLink({ label, icon, path, onClick }: NavBarLinkProps) {
  const theme = useMantineTheme();
  const location = useLocation();

  if (path) {
    return (
      <NavLink
        label={label}
        component={Link}
        to={path}
        active={location.pathname === path}
        sx={{ borderRadius: theme.radius.sm }}
        icon={icon}
        onClick={onClick}
      />
    );
  }

  return (
    <NavLink
      label={label}
      sx={{ borderRadius: theme.radius.sm }}
      icon={icon}
      onClick={onClick}
    />
  );
}
