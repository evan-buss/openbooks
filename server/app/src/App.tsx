import {
  AppShell,
  Burger,
  ColorScheme,
  ColorSchemeProvider,
  createEmotionCache,
  MantineProvider,
  MediaQuery
} from "@mantine/core";
import { useLocalStorage } from "@mantine/hooks";
import { NotificationsProvider } from "@mantine/notifications";
import NotificationDrawer from "./components/SideBar/NotificationDrawer";
import SidebarNeo from "./components/SidebarNeo/Sidebar";
import SearchPage from "./pages/SearchPage";

const emotionCache = createEmotionCache({ key: "openbooks" });

export default function App() {
  const [colorScheme, setColorScheme] = useLocalStorage({
    key: "color-scheme",
    defaultValue: "light" as ColorScheme,
    getInitialValueInEffect: true
  });

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
              "#1a53ff",
              "#0039e6",
              "#002db4",
              "#002082",
              "#001351",
              "#000621"
            ]
          }
        }}>
        <NotificationsProvider position="top-center">
          <AppShell
            navbar={<SidebarNeo />}
            padding={0}
            styles={(theme) => ({
              main: {
                backgroundColor:
                  theme.colorScheme === "dark"
                    ? theme.colors.dark[8]
                    : theme.colors.gray[0]
              }
            })}>
            <div className="flex flex-row max-h-screen min-h-screen flex-nowrap">
              <MediaQuery largerThan="sm" styles={{ display: "none" }}>
                <Burger opened={false} size="sm" mr="xl" />
              </MediaQuery>
              <SearchPage />
              <NotificationDrawer />
            </div>
          </AppShell>
        </NotificationsProvider>
      </MantineProvider>
    </ColorSchemeProvider>
  );
}
