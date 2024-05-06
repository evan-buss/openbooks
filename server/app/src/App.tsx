import {
  AppShell,
  ColorSchemeScript,
  createTheme,
  MantineProvider,
  px
} from "@mantine/core";
import { Notifications } from "@mantine/notifications";
import { RootState, useAppSelector } from "./state/store";
// import "@mantine/core/styles.css";
import "@mantine/core/styles/global.css";

import "@mantine/core/styles/UnstyledButton.css";
import "@mantine/core/styles/ScrollArea.css";
import "@mantine/core/styles/VisuallyHidden.css";
import "@mantine/core/styles/Paper.css";
import "@mantine/core/styles/Popover.css";
import "@mantine/core/styles/CloseButton.css";
import "@mantine/core/styles/Group.css";
import "@mantine/core/styles/Loader.css";
import "@mantine/core/styles/Overlay.css";
import "@mantine/core/styles/ModalBase.css";
import "@mantine/core/styles/Input.css";
import "@mantine/core/styles/Flex.css";

import "@mantine/core/styles/Button.css";

import "@mantine/core/styles/ActionIcon.css";
import "@mantine/core/styles/AppShell.css";
import "@mantine/core/styles/Badge.css";
import "@mantine/core/styles/Center.css";
import "@mantine/core/styles/Code.css";
import "@mantine/core/styles/Divider.css";
import "@mantine/core/styles/Image.css";
import "@mantine/core/styles/Indicator.css";
import "@mantine/core/styles/Menu.css";
import "@mantine/core/styles/Modal.css";
import "@mantine/core/styles/Notification.css";
import "@mantine/core/styles/Radio.css";
import "@mantine/core/styles/SegmentedControl.css";
import "@mantine/core/styles/Stack.css";
import "@mantine/core/styles/Switch.css";
import "@mantine/core/styles/Table.css";
import "@mantine/core/styles/Text.css";
import "@mantine/core/styles/Tooltip.css";
import "@mantine/core/styles/TypographyStylesProvider.css";

import "@mantine/notifications/styles.css";

import Sidebar from "./components/sidebar/Sidebar";
import SearchPage from "./pages/SearchPage";
import { NotificationDrawer } from "./components/drawer/NotificationDrawer";
import classes from "./App.module.css";
import { useSelector } from "react-redux";

const theme = createTheme({
  activeClassName: "",
  primaryColor: "blue",
  primaryShade: { light: 4, dark: 2 },
  colors: {
    blue: [
      "#e0ecff",
      "#b0c6ff",
      "#7e9fff",
      "#4c79ff",
      "#3366ff",
      "#0039e6",
      "#002db4",
      "#002082",
      "#001351",
      "#000621"
    ]
  },
  components: {
    ActionIcon: {
      defaultProps: {
        radius: "md",
        color: "blue",
        variant: "subtle"
      }
    }
  }
});

export default function App() {
  const opened = useAppSelector((state) => state.state.isSidebarOpen);
  const notificationsOpen = useSelector(
    (state: RootState) => state.notifications.isOpen
  );

  return (
    <>
      <ColorSchemeScript defaultColorScheme="auto" />
      <MantineProvider defaultColorScheme="auto" theme={theme}>
        <Notifications position="top-center" />
        <AppShell
          navbar={{
            width: 300,
            breakpoint: "sm",
            collapsed: { mobile: !opened, desktop: !opened }
          }}
          padding={px("1.5rem")}
          aside={{
            width: 300,
            breakpoint: "",
            collapsed: {
              mobile: !notificationsOpen,
              desktop: !notificationsOpen
            }
          }}>
          <Sidebar />
          <AppShell.Main className={classes.main}>
            <SearchPage />
          </AppShell.Main>
          <AppShell.Aside>
            <NotificationDrawer />
          </AppShell.Aside>
        </AppShell>
      </MantineProvider>
    </>
  );
}
