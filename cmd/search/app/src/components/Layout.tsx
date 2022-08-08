import { AppShell, useMantineTheme } from "@mantine/core";
import { useState } from "react";
import { Outlet } from "react-router-dom";
import { NavBar } from "./NavBar";
import { Sidebar } from "./Sidebar";

export function Layout() {
  const theme = useMantineTheme();
  const [opened, setOpened] = useState(false);
  return (
    <AppShell
      styles={{
        main: {
          background:
            theme.colorScheme === "dark"
              ? theme.colors.dark[8]
              : theme.colors.gray[0],
        },
      }}
      navbarOffsetBreakpoint="sm"
      asideOffsetBreakpoint="sm"
      navbar={<Sidebar opened={opened} setOpened={setOpened} />}
      header={<NavBar opened={opened} setOpened={setOpened} />}
    >
      <Outlet />
    </AppShell>
  );
}
