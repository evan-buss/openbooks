import {
  Badge,
  Menu,
  Popover,
  Position,
  Spinner,
  Text,
  Tooltip
} from "evergreen-ui";
import { Book, Download, Trash } from "phosphor-react";
import React from "react";
import { useDeleteBookMutation, useGetBooksQuery } from "../../state/api";
import { downloadFile } from "../../state/util";

export default function BookList() {
  const { data, isLoading, isSuccess, isError } = useGetBooksQuery(null);
  const [deleteBook] = useDeleteBookMutation();

  if (isLoading) {
    return (
      <div className="flex justify-center m-4">
        <Spinner />
      </div>
    );
  }

  if (isError) {
    return (
      <p className="text-center my-4 text-gray-500 text-sm">
        Book persistence disabled.
      </p>
    );
  }

  if (isSuccess && data?.length === 0) {
    return (
      <p className="text-center my-4 text-gray-500 text-sm">
        No previous downloads.
      </p>
    );
  }

  return (
    <>
      {data?.map((book) => (
        <Popover
          key={book.time}
          position={Position.BOTTOM}
          content={
            <Menu>
              <Menu.Item
                icon={<Download color="black" size={18} weight="bold" />}
                onClick={() => downloadFile(book.downloadLink)}>
                Download
              </Menu.Item>
              <Menu.Item
                icon={<Trash color="#d45050" size={18} weight="bold" />}
                intent="danger"
                onClick={() => deleteBook(book.name)}>
                Delete
              </Menu.Item>
            </Menu>
          }>
          <div className="border cursor-pointer p-2 m-1 rounded hover:bg-hover-blue flex justify-between items-center">
            <Book size={18} />
            <Tooltip
              position={Position.RIGHT}
              showDelay={500}
              hideDelay={0}
              content={book.name}>
              <Text
                width={140}
                marginLeft={24}
                flex={1}
                whiteSpace="nowrap"
                display="block"
                fontWeight={500}
                overflow="hidden"
                textOverflow="ellipsis">
                {book.name}
              </Text>
            </Tooltip>
            <Badge color="blue">
              {new Date(book.time).toLocaleDateString("en-US")}
            </Badge>
          </div>
        </Popover>
      ))}
    </>
  );
}
