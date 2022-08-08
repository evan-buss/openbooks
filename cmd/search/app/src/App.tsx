import {
  ColorScheme,
  ColorSchemeProvider,
  MantineProvider,
} from "@mantine/core";
import { useColorScheme } from "@mantine/hooks";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { AuthRoute } from "./components/AuthRoute";
import { Layout } from "./components/Layout";
import DashboardPage from "./pages/dashboard";
import { LogInPage } from "./pages/login";
import { SettingsPage } from "./pages/settings";
import { AuthProvider } from "./util/auth";

const queryClient = new QueryClient();

export default function App() {
  // hook will return either 'dark' or 'light' on client
  // and always 'light' during ssr as window.matchMedia is not available
  const preferredColorScheme = useColorScheme("dark");
  const [colorScheme, setColorScheme] =
    useState<ColorScheme>(preferredColorScheme);
  const toggleColorScheme = (value?: ColorScheme) =>
    setColorScheme(value || (colorScheme === "dark" ? "light" : "dark"));

  return (
    <QueryClientProvider client={queryClient}>
      <ColorSchemeProvider
        colorScheme={colorScheme}
        toggleColorScheme={toggleColorScheme}
      >
        <MantineProvider
          theme={{ colorScheme, defaultRadius: "md" }}
          withGlobalStyles
          withNormalizeCSS
        >
          <BrowserRouter>
            <AuthProvider>
              <Routes>
                <Route
                  path="/"
                  element={
                    <AuthRoute>
                      <Layout />
                    </AuthRoute>
                  }
                >
                  <Route index element={<DashboardPage />} />
                  <Route path="/settings" element={<SettingsPage />} />
                </Route>
                <Route path="/login" element={<LogInPage />} />
              </Routes>
            </AuthProvider>
          </BrowserRouter>
        </MantineProvider>
      </ColorSchemeProvider>
    </QueryClientProvider>
  );
}
