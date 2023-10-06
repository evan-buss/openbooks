import {
  ActionIcon,
  AppShell,
  Group,
  SegmentedControl,
  Text,
  Tooltip,
  useMantineColorScheme
} from "@mantine/core";
import { useLocalStorage } from "@mantine/hooks";
import {
  BellSimple,
  MoonStars,
  Sidebar as SidebarIcon,
  Sun,
  X
} from "@phosphor-icons/react";
import { toggleDrawer } from "../../state/notificationSlice";
import { toggleSidebar } from "../../state/stateSlice";
import { useAppDispatch, useAppSelector } from "../../state/store";
import History from "./History";
import Library from "./Library";
import classes from "./Sidebar.module.css";
import { ServerMenu } from "./ServerMenu";

export default function Sidebar() {
  const { colorScheme, toggleColorScheme } = useMantineColorScheme();

  const dispatch = useAppDispatch();
  const connected = useAppSelector((store) => store.state.isConnected);
  const username = useAppSelector((store) => store.state.username);
  const serverName = useAppSelector(
    (store) => store.connection.selectedServer.name
  );

  const [index, setIndex] = useLocalStorage<"books" | "history">({
    key: "sidebar-state",
    defaultValue: "history"
  });

  return (
    <AppShell.Navbar className={classes.navbar}>
      {/*<AppShell.Navbar>*/}
      <AppShell.Section p="sm">
        <Group justify="space-between">
          <Text fw="bold" size="lg">
            OpenBooks
          </Text>
          <Group>
            <Tooltip
              label={`OpenBooks server ${
                connected ? "connected" : "disconnected"
              }.`}>
              <ActionIcon
                disabled={!connected}
                onClick={() => dispatch(toggleDrawer())}>
                <BellSimple weight="bold" size={18} />
              </ActionIcon>
            </Tooltip>
            <ActionIcon
              hiddenFrom="sm"
              variant="subtle"
              color="dark"
              onClick={() => dispatch(toggleSidebar())}>
              <X weight="bold" size={18} />
            </ActionIcon>
          </Group>
        </Group>

        <Text size="sm" c="dimmed">
          Download eBooks from {serverName}
        </Text>

        <SegmentedControl
          size="sm"
          styles={(theme) => ({
            root: {
              marginTop: theme.spacing.md
            },
            label: {
              fontSize: theme.fontSizes.xs
            }
          })}
          value={index}
          onChange={(value: "books" | "history") => setIndex(value)}
          data={[
            { label: "Search History", value: "history" },
            { label: "Previous Downloads", value: "books" }
          ]}
          fullWidth
        />
      </AppShell.Section>

      <AppShell.Section grow p="xs" style={{ overflow: "auto" }}>
        {index === "history" ? <History /> : <Library />}
      </AppShell.Section>

      <AppShell.Section p="sm" className={classes.footer}>
        <Group justify="space-between" wrap="nowrap">
          <Group>
            <ServerMenu connected={!!username} />
            {username ? (
              <Text
                size="sm"
                lineClamp={1}
                style={{
                  maxWidth: 150,
                  overflow: "hidden",
                  textOverflow: "ellipsis",
                  whiteSpace: "nowrap"
                }}>
                {username}
              </Text>
            ) : (
              <Text size="sm">Not connected.</Text>
            )}
          </Group>

          <Group align="end" gap="xs">
            <ActionIcon onClick={() => toggleColorScheme()}>
              {colorScheme === "dark" ? (
                <Sun size={18} weight="bold" />
              ) : (
                <MoonStars size={18} weight="bold" />
              )}
            </ActionIcon>
            <ActionIcon onClick={() => dispatch(toggleSidebar())}>
              <SidebarIcon weight="bold" size={18} />
            </ActionIcon>
          </Group>
        </Group>
      </AppShell.Section>
    </AppShell.Navbar>
  );
}
