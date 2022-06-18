import React, { useState } from "react";
import { ParseError } from "../../state/messages";
import { majorScale, Table, Text } from "evergreen-ui";
import { useEvent } from "react-use";
import debounce from "lodash/debounce";
import { useHeight } from "../../state/hooks";

export interface Props {
  errors?: ParseError[];
  setSearchQuery: (query: string) => void;
}

export default function ErrorsGrid({ errors, setSearchQuery }: Props) {
  const [lineFilter, setLineFilter] = useState<string>("");
  const [errorFilter, setErrorFilter] = useState<string>("");
  const height = useHeight();

  const onSelectionChange = () => {
    const selection = document.getSelection()?.toString();
    if (selection) setSearchQuery(selection);
  };

  useEvent("selectionchange", debounce(onSelectionChange, 100), document);

  if (errors?.length === 0) {
    return (
      <Text textAlign="center" padding={majorScale(3)}>
        No results for search.
      </Text>
    );
  }

  return (
    <>
      <Text className="w-full mb-2 text-gray-500" size={400}>
        These results could not be parsed to due to their non-standard format.
        To download, copy the line up to the <code>::INFO::</code> or file size
        at the end and paste into the text box above.
      </Text>
      <Table
        key={height}
        flex={1}
        display="flex"
        flexDirection="column"
        width="100%"
        className="rounded-lg">
        <Table.Head background="white" className="rounded-t-md">
          <Table.SearchHeaderCell
            onChange={(e) => setLineFilter(e)}
            placeholder="LINE"
            flexBasis={250}
            flexGrow={1}
            flexShrink={0}
          />
          <Table.SearchHeaderCell
            onChange={(e) => setErrorFilter(e)}
            placeholder="PARSE ERROR"
            flexBasis={350}
            flexGrow={0}
            flexShrink={0}
          />
        </Table.Head>
        <Table.VirtualBody>
          {errors
            ?.filter(
              (x) =>
                x.line.toLowerCase().includes(lineFilter.toLowerCase()) &&
                x.error?.toLowerCase().includes(errorFilter.toLowerCase())
            )
            .map((book, i) => (
              <Table.Row key={book.line + i} isSelectable>
                <Table.TextCell
                  flexBasis={250}
                  flexGrow={1}
                  flexShrink={0}
                  cursor="text">
                  <code>{book.line}</code>
                </Table.TextCell>
                <Table.TextCell flexBasis={350} flexGrow={0} flexShrink={0}>
                  {book.error}
                </Table.TextCell>
              </Table.Row>
            ))}
        </Table.VirtualBody>
      </Table>
    </>
  );
}
