import {
  Button,
  Center,
  createStyles,
  Group,
  Image,
  MediaQuery,
  Stack,
  TextInput
} from "@mantine/core";
import { MagnifyingGlass, Warning } from "phosphor-react";
import { FormEvent, useEffect, useMemo, useState } from "react";
import { useSelector } from "react-redux";
import image from "../assets/reading.svg";
import { BooksGrid } from "../components/BooksGrid/BooksGrid";
import ErrorsGrid from "../components/ErrorsGrid/ErrorsGrid";
import BookTable from "../components/Tables/BookTable";
import ErrorTable from "../components/Tables/ErrorTable";
import { MessageType } from "../state/messages";
import { sendMessage, sendSearch } from "../state/stateSlice";
import { RootState, useAppDispatch } from "../state/store";

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

export default function SearchPage({ legacy }: { legacy?: boolean }) {
  const dispatch = useAppDispatch();
  const activeItem = useSelector((store: RootState) => store.state.activeItem);
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

  const renderBody = () => {
    if (activeItem === null) {
      return (
        <Center style={{ height: "100%", width: "100%" }}>
          <MediaQuery smallerThan="md" styles={{ display: "none" }}>
            <Image width={600} fit="contain" src={image} alt="person reading" />
          </MediaQuery>
          <MediaQuery largerThan="md" styles={{ display: "none" }}>
            <Image width={300} fit="contain" src={image} alt="person reading" />
          </MediaQuery>
        </Center>
      );
    }

    if (legacy) {
      return errorMode ? (
        <ErrorsGrid
          errors={activeItem?.errors ?? []}
          setSearchQuery={setSearchQuery}
        />
      ) : (
        <BooksGrid books={activeItem?.results ?? []} />
      );
    }

    return errorMode ? errorTable : bookTable;
  };

  return (
    <Stack
      spacing={0}
      align="center"
      sx={(theme) => ({ width: "100%", margin: theme.spacing.xl })}>
      <form className={classes.wFull} onSubmit={(e) => searchHandler(e)}>
        <Group
          noWrap
          spacing="xl"
          sx={(theme) => ({ marginBottom: theme.spacing.xl })}>
          <TextInput
            className={classes.wFull}
            disabled={activeItem !== null && !activeItem.results}
            value={searchQuery}
            onChange={(e: any) => setSearchQuery(e.target.value)}
            placeholder={
              errorMode ? "Download a book manually." : "Search for a book."
            }
            radius="md"
            size="md"
            type="search"
            icon={<MagnifyingGlass weight="bold" size={22} />}
            required
          />

          <Button
            type="submit"
            color={theme.colorScheme === "dark" ? "brand.2" : "brand"}
            disabled={!validInput}
            size="md"
            radius="md"
            variant="outline">
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
      {renderBody()}
    </Stack>
  );
}
