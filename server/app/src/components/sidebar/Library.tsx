import {
  Badge,
  Button,
  Center,
  Loader,
  Menu,
  Stack,
  Text,
  Tooltip
} from "@mantine/core";
import { Book as BookIcon, Download, Trash } from "@phosphor-icons/react";
import { Book, useDeleteBookMutation, useGetBooksQuery } from "../../state/api";
import { downloadFile } from "../../state/util";
import classes from "./SidebarButton.module.css";
import { useAutoAnimate } from "@formkit/auto-animate/react";

export default function Library() {
  const { data, isLoading, isSuccess, isError } = useGetBooksQuery();
  const [parent] = useAutoAnimate(/* optional config */);

  if (isLoading && !data) {
    return (
      <Center>
        <Loader />
      </Center>
    );
  }

  if (isError) {
    return (
      <Center>
        <Text c="dimmed" size="sm">
          Book persistence disabled.
        </Text>
      </Center>
    );
  }

  if (isSuccess && data?.length === 0) {
    return (
      <Center>
        <Text c="dimmed" size="sm">
          No previous downloads.
        </Text>
      </Center>
    );
  }

  return (
    <Stack gap="xs" ref={parent}>
      {data?.map((book) => <LibraryCard key={book.name} book={book} />)}
    </Stack>
  );
}

interface LibraryCardProps {
  book: Book;
}

function LibraryCard({ book }: LibraryCardProps) {
  const [deleteBook] = useDeleteBookMutation();

  return (
    <Menu shadow="md">
      <Menu.Target>
        <Tooltip label={book.name} openDelay={1_000}>
          <Button
            classNames={{
              root: classes.root,
              section: classes.section,
              label: classes.label
            }}
            key={book.name}
            radius="sm"
            variant="outline"
            fullWidth
            leftSection={<BookIcon weight="bold" size={18} />}
            rightSection={
              <Badge color="brand" radius="sm" size="sm" variant="light">
                {new Date(book.time).toLocaleDateString("en-US")}
              </Badge>
            }>
            {book.name}
          </Button>
        </Tooltip>
      </Menu.Target>

      <Menu.Dropdown>
        <Menu.Item
          leftSection={<Download weight="bold" />}
          onClick={() => downloadFile(book.downloadLink)}>
          Download
        </Menu.Item>

        <Menu.Item
          color="red"
          leftSection={<Trash size={18} weight="bold" />}
          onClick={() => deleteBook(book.name)}>
          Delete
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  );
}
