import {
  ActionIcon,
  Burger,
  Header,
  MediaQuery,
  Text,
  useMantineColorScheme,
  useMantineTheme,
} from "@mantine/core";
import { CloudMoon, Sun } from "phosphor-react";
import { Dispatch, SetStateAction } from "react";

export interface NavBarProps {
  opened: boolean;
  setOpened: Dispatch<SetStateAction<boolean>>;
}

export function NavBar({ opened, setOpened }: NavBarProps) {
  const theme = useMantineTheme();
  const { colorScheme, toggleColorScheme } = useMantineColorScheme();
  const dark = colorScheme === "dark";

  return (
    <Header height={60} p="md">
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          height: "100%",
        }}
      >
        <MediaQuery largerThan="sm" styles={{ display: "none" }}>
          <Burger
            opened={opened}
            onClick={() => setOpened((o) => !o)}
            size="sm"
            color={theme.colors.gray[6]}
            mr="xl"
          />
        </MediaQuery>

        <Text size="lg">OpenBooks Search</Text>

        <ActionIcon
          color={dark ? "yellow" : "blue"}
          onClick={() => toggleColorScheme()}
          title="Toggle color scheme"
        >
          {dark ? <Sun size={16} /> : <CloudMoon size={16} />}
        </ActionIcon>
      </div>
    </Header>
  );
}
