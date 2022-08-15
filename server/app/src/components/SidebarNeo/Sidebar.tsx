import {
  ActionIcon,
  Burger,
  createStyles,
  Group,
  MediaQuery,
  Navbar,
  ScrollArea,
  SegmentedControl,
  Text,
  Transition,
  useMantineColorScheme
} from "@mantine/core";
import { IconMoonStars, IconSun } from "@tabler/icons";
import { IdentificationBadge } from "phosphor-react";
import { useState } from "react";
import { useSelector } from "react-redux";
import { useLocalStorage } from "react-use";
import { toggleDrawer } from "../../state/notificationSlice";
import { RootState, useAppDispatch } from "../../state/store";
import BookList from "../SideBar/BookList";
import Pulse from "../SideBar/Pulse";
import SearchHistoryNeo from "./History";

const useStyles = createStyles((theme, _params, getRef) => {
  return {
    navbar: {
      backgroundColor:
        theme.colorScheme === "dark" ? theme.colors.dark[7] : theme.white
    },
    footer: {
      borderTop: `1px solid ${
        theme.colorScheme === "dark"
          ? theme.colors.dark[4]
          : theme.colors.gray[3]
      }`,
      paddingTop: theme.spacing.sm
    }
  };
});

export default function SidebarNeo() {
  const { classes } = useStyles();
  const [opened, setOpened] = useState(true);
  const username = useSelector((store: RootState) => store.state.username);
  const [index, setIndex] = useLocalStorage<"books" | "history">(
    "newIndex",
    "history"
  );

  const { colorScheme, toggleColorScheme } = useMantineColorScheme();
  const connected = useSelector((store: RootState) => store.state.isConnected);
  const dispatch = useAppDispatch();

  return (
    <Navbar
      width={{ sm: 300 }}
      hiddenBreakpoint="sm"
      hidden={!opened}
      className={classes.navbar}>
      <Navbar.Section p="sm">
        <Group position="apart">
          <Text weight="bold" size="lg">
            OpenBooks
          </Text>
          <div>
            <MediaQuery largerThan="sm" styles={{ display: "none" }}>
              <Burger
                opened={opened}
                onClick={() => setOpened((o) => !o)}
                size="sm"
                mr="xl"
              />
            </MediaQuery>
            <Pulse
              enabled={connected}
              onClick={() => dispatch(toggleDrawer())}
            />
          </div>
        </Group>

        <Text size="sm" color="dimmed">
          Download eBooks from IRC Highway
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
          transitionDuration={500}
          transitionTimingFunction="ease"
          value={index}
          onChange={(value: "books" | "history") => setIndex(value)}
          data={[
            { label: "Search History", value: "history" },
            { label: "Previous Downloads", value: "books" }
          ]}
          fullWidth
        />
      </Navbar.Section>

      <Navbar.Section grow component={ScrollArea} p="xs">
        <Transition
          mounted={index === "history"}
          transition="fade"
          duration={1000}
          exitDuration={0}
          timingFunction="ease">
          {(styles) => (
            <div style={styles}>
              <SearchHistoryNeo />
            </div>
          )}
        </Transition>

        <Transition
          mounted={index === "books"}
          transition="fade"
          duration={1000}
          exitDuration={0}
          timingFunction="ease">
          {(styles) => (
            <div style={styles}>
              <BookList />
            </div>
          )}
        </Transition>
      </Navbar.Section>

      <Navbar.Section className={classes.footer} p="sm">
        <Group position="apart">
          {username && (
            <Group>
              <IdentificationBadge size={24} className="mr-4" />
              <Text size="sm">{username}</Text>
            </Group>
          )}

          <ActionIcon
            onClick={() => toggleColorScheme()}
            size="lg"
            sx={(theme) => ({
              backgroundColor:
                theme.colorScheme === "dark"
                  ? theme.colors.dark[6]
                  : theme.colors.gray[0],
              color:
                theme.colorScheme === "dark"
                  ? theme.colors.yellow[4]
                  : theme.colors.brand[4]
            })}>
            {colorScheme === "dark" ? (
              <IconSun size={18} />
            ) : (
              <IconMoonStars size={18} />
            )}
          </ActionIcon>
        </Group>
      </Navbar.Section>
    </Navbar>
  );
}
