import {
  ActionIcon,
  Button,
  Center,
  createStyles,
  Group,
  Image,
  MediaQuery,
  Stack,
  TextInput,
  Title
} from "@mantine/core";
import { MagnifyingGlass, Sidebar, Warning } from "phosphor-react";
import { FormEvent, useEffect, useMemo, useState } from "react";
import image from "../assets/reading.svg";
import BookTable from "../components/tables/BookTable";
import ErrorTable from "../components/tables/ErrorTable";
import { MessageType } from "../state/messages";
import { sendMessage, sendSearch, toggleSidebar } from "../state/stateSlice";
import { useAppDispatch, useAppSelector } from "../state/store";

const useStyles = createStyles(
  (theme, { errorMode }: { errorMode: boolean }) => ({
    stack: {
      minWidth: "100%",
      margin: theme.spacing.xl,
      backgroundColor: theme.colors.blue[0]
    },
    wFull: {
      width: "100%"
    },
    errorToggle: {
      "alignSelf": "start",
      "height": "24px",
      "marginBottom": theme.spacing.xs,
      "fontWeight": 500,
      "color":
        theme.colorScheme === "dark"
          ? errorMode
            ? theme.colors.dark[8]
            : theme.colors.dark[2]
          : errorMode
          ? theme.colors.white
          : theme.colors.dark[3],
      "&:hover": {
        backgroundColor:
          theme.colorScheme === "dark"
            ? errorMode
              ? theme.colors.brand[3]
              : theme.colors.dark[7]
            : errorMode
            ? theme.colors.brand[5]
            : theme.colors.gray[1]
      }
    }
  })
);

export default function SearchPage() {
  const dispatch = useAppDispatch();
  const activeItem = useAppSelector((store) => store.state.activeItem);
  const opened = useAppSelector((store) => store.state.isSidebarOpen);

  const [searchQuery, setSearchQuery] = useState("");
  const [showErrors, setShowErrors] = useState(false);

  const hasErrors = (activeItem?.errors ?? []).length > 0;
  const errorMode = showErrors && activeItem;
  const validInput = errorMode
    ? searchQuery.startsWith("!")
    : searchQuery !== "";

  const { classes, theme } = useStyles({ errorMode: !!errorMode });

  useEffect(() => {
    setShowErrors(false);
  }, [activeItem]);

  const searchHandler = (event: FormEvent) => {
    event.preventDefault();

    if (errorMode) {
      dispatch(
        sendMessage({
          type: MessageType.DOWNLOAD,
          payload: { book: searchQuery }
        })
      );
    } else {
      dispatch(sendSearch(searchQuery));
    }

    setSearchQuery("");
  };

  const bookTable = useMemo(
    () => <BookTable books={activeItem?.results ?? []} />,
    [activeItem?.results]
  );

  const errorTable = useMemo(
    () => (
      <ErrorTable
        errors={activeItem?.errors ?? []}
        setSearchQuery={setSearchQuery}
      />
    ),
    [activeItem?.errors]
  );

  return (
    <Stack
      spacing={0}
      align="center"
      sx={(theme) => ({ width: "100%", margin: theme.spacing.xl })}>
      <form className={classes.wFull} onSubmit={(e) => searchHandler(e)}>
        <Group
          noWrap
          spacing="md"
          sx={(theme) => ({ marginBottom: theme.spacing.md })}>
          {!opened && (
            <ActionIcon size="lg" onClick={() => dispatch(toggleSidebar())}>
              <Sidebar weight="bold" size={20}></Sidebar>
            </ActionIcon>
          )}
          <TextInput
            className={classes.wFull}
            variant="filled"
            disabled={activeItem !== null && !activeItem.results}
            value={searchQuery}
            onChange={(e: any) => setSearchQuery(e.target.value)}
            placeholder={
              errorMode ? "Download a book manually." : "Search for a book."
            }
            radius="md"
            type="search"
            icon={<MagnifyingGlass weight="bold" size={22} />}
            required
          />

          <Button
            type="submit"
            color={theme.colorScheme === "dark" ? "brand.2" : "brand"}
            disabled={!validInput}
            radius="md"
            variant={validInput ? "gradient" : "default"}
            gradient={{ from: "brand.4", to: "brand.3" }}>
            {errorMode ? "Download" : "Search"}
          </Button>
        </Group>
      </form>
      {hasErrors && (
        <Button
          className={classes.errorToggle}
          variant={errorMode ? "filled" : "subtle"}
          onClick={() => setShowErrors((show) => !show)}
          leftIcon={<Warning size={18} />}
          size="xs">
          {activeItem?.errors?.length} Parsing{" "}
          {activeItem?.errors?.length === 1 ? "Error" : "Errors"}
        </Button>
      )}
      {!activeItem ? (
        <Center style={{ height: "100%", width: "100%" }}>
          <Stack align="center">
            <Title weight="normal" align="center">
              Search a book to get started.
            </Title>
            <MediaQuery smallerThan="md" styles={{ display: "none" }}>
              <Image
                width={600}
                fit="contain"
                src={image}
                alt="person reading"
              />
            </MediaQuery>
            <MediaQuery largerThan="md" styles={{ display: "none" }}>
              <Image
                width={300}
                fit="contain"
                src={image}
                alt="person reading"
              />
            </MediaQuery>
          </Stack>
        </Center>
      ) : errorMode ? (
        errorTable
      ) : (
        bookTable
      )}
    </Stack>
  );
}
