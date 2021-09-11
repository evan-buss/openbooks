import {
  Text,
  BookIcon,
  Badge,
  Spinner,
  Popover,
  Position,
  TrashIcon,
  Menu,
  DownloadIcon
} from "evergreen-ui";
import React from "react";
import { useDeleteBookMutation, useGetBooksQuery } from "../../state/api";
import { downloadFile } from "../../state/util";

export default function BookList() {
  const { data, isLoading } = useGetBooksQuery(null);
  const [deleteBook] = useDeleteBookMutation();

  if (isLoading) {
    return <Spinner></Spinner>;
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
                icon={DownloadIcon}
                intent="info"
                onClick={() => downloadFile(book.downloadLink)}>
                Download
              </Menu.Item>
              <Menu.Item
                icon={TrashIcon}
                intent="danger"
                onClick={() => deleteBook(book.name)}>
                Delete
              </Menu.Item>
            </Menu>
          }>
          <div className="border cursor-pointer p-2 m-1 shadow flex justify-between items-center">
            <BookIcon size={15} color="#234361" />
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
            <Badge color="green">
              {new Date(book.time).toLocaleDateString("en-US")}
            </Badge>
          </div>
        </Popover>
      ))}
    </>
  );
}
