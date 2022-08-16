import {
  AppShell,
  Burger,
  ColorScheme,
  ColorSchemeProvider,
  MantineProvider,
  MediaQuery
} from "@mantine/core";
import { useState } from "react";
import NotificationDrawer from "./components/SideBar/NotificationDrawer";
import SidebarNeo from "./components/SidebarNeo/Sidebar";
import SearchPage from "./pages/SearchPage";

function App() {
  const [colorScheme, setColorScheme] = useState<ColorScheme>("light");

  const toggleColorScheme = (value?: ColorScheme) =>
    setColorScheme(value || (colorScheme === "dark" ? "light" : "dark"));

  return (
    <ColorSchemeProvider
      colorScheme={colorScheme}
      toggleColorScheme={toggleColorScheme}>
      <MantineProvider
        withGlobalStyles
        withNormalizeCSS
        theme={{
          colorScheme: colorScheme,
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
            {/* <Sidebar /> */}
            <MediaQuery largerThan="sm" styles={{ display: "none" }}>
              <Burger opened={false} size="sm" mr="xl" />
            </MediaQuery>
            <SearchPage />
            <NotificationDrawer />
          </div>
        </AppShell>
      </MantineProvider>
    </ColorSchemeProvider>
  );
}

export default App;
