import {
  AppShell,
  ColorScheme,
  ColorSchemeProvider,
  createEmotionCache,
  createStyles,
  MantineProvider
} from "@mantine/core";
import { useLocalStorage } from "@mantine/hooks";
import { Notifications } from "@mantine/notifications";
import NotificationDrawer from "./components/drawer/NotificationDrawer";
import Sidebar from "./components/sidebar/Sidebar";
import SearchPage from "./pages/SearchPage";
import { useAppDispatch, useAppSelector } from "./state/store";

const emotionCache = createEmotionCache({ key: "openbooks" });

const useStyles = createStyles(() => ({
  burger: {
    position: "absolute",
    bottom: 0,
    left: 0
  },
  wrapper: {
    boxSizing: "border-box",
    display: "flex",
    flexWrap: "nowrap",
    maxHeight: "100vh",
    minHeight: "100vh"
  }
}));

export default function App() {
  const { classes } = useStyles();
  const [colorScheme, setColorScheme] = useLocalStorage({
    key: "color-scheme",
    defaultValue: "light" as ColorScheme,
    getInitialValueInEffect: true
  });

  const dispatch = useAppDispatch();
  const open = useAppSelector((state) => state.state.isSidebarOpen);

  return (
    <ColorSchemeProvider
      colorScheme={colorScheme}
      toggleColorScheme={() =>
        setColorScheme((color) => (color === "dark" ? "light" : "dark"))
      }>
      <MantineProvider
        emotionCache={emotionCache}
        withGlobalStyles
        withNormalizeCSS
        theme={{
          colorScheme,
          activeStyles: { transform: "none" },
          primaryColor: "brand",
          primaryShade: { light: 4, dark: 2 },
          colors: {
            brand: [
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
                color: "brand"
              }
            }
          }
        }}>
        <Notifications position="top-center" />
        <AppShell
          aside={<Sidebar />}
          padding={0}
          styles={(theme) => ({
            main: {
              backgroundColor:
                theme.colorScheme === "dark"
                  ? theme.colors.dark[8]
                  : theme.colors.gray[0]
            }
          })}>
          <div className={classes.wrapper}>
            <SearchPage />
            <NotificationDrawer />
          </div>
        </AppShell>
      </MantineProvider>
    </ColorSchemeProvider>
  );
}
