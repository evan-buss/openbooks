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
import { Book, Download, Trash } from "phosphor-react";
import { useDeleteBookMutation, useGetBooksQuery } from "../../state/api";
import { downloadFile } from "../../state/util";
import { useSidebarButtonStyle } from "./styles";

export default function Library() {
  const { data, isLoading, isSuccess, isError } = useGetBooksQuery(null);
  const [deleteBook] = useDeleteBookMutation();
  const { classes } = useSidebarButtonStyle({});

  if (isLoading) {
    return (
      <Center>
        <Loader />
      </Center>
    );
  }

  if (isError) {
    return (
      <Center>
        <Text color="dimmed" size="sm">
          Book persistence disabled.
        </Text>
      </Center>
    );
  }

  if (isSuccess && data?.length === 0) {
    return (
      <Center>
        <Text color="dimmed" size="sm">
          No previous downloads.
        </Text>
      </Center>
    );
  }

  return (
    <Stack spacing="xs" sx={{ width: "278px" }}>
      {new Array(20).fill(data![0])?.map((book) => (
        <Menu shadow="md" width={200}>
          <Menu.Target>
            <Tooltip label={book.name} openDelay={1_000}>
              <Button
                classNames={classes}
                radius="sm"
                variant="outline"
                leftIcon={<Book weight="bold" />}
                rightIcon={
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
              icon={<Download weight="bold" />}
              onClick={() => downloadFile(book.downloadLink)}>
              Download
            </Menu.Item>

            <Menu.Item
              color="red"
              icon={<Trash size={18} weight="bold" />}
              onClick={() => deleteBook(book.name)}>
              Delete
            </Menu.Item>
          </Menu.Dropdown>
        </Menu>
      ))}
    </Stack>
  );
}
