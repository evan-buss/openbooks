import {
  Option,
  SelectMenu,
  SelectMenuItem,
  Table,
  Image,
  TextDropdownButton,
  StatusIndicator
} from "evergreen-ui";
import React from "react";

interface Props {
  flexBasis?: number;
  columnTitle: string;
  menuTitle: string;
  options: string[];
  selected: string[];
  setSelected: React.Dispatch<React.SetStateAction<string[]>>;
  itemRenderer?: any;
}

const SelectMenuHeader = (props: Props) => {
  const menuOptions: SelectMenuItem[] = props.options.map((opt) => ({
    label: opt,
    value: opt
  }));

  return (
    <Table.TextHeaderCell
      flexBasis={props.flexBasis ?? 100}
      flexGrow={0}
      flexShrink={0}>
      <SelectMenu
        isMultiSelect
        title={props.menuTitle}
        options={menuOptions}
        selected={props.selected}
        itemRenderer={props.itemRenderer}
        onSelect={(item) =>
          props.setSelected((curr) => [...curr, item.value.toString()])
        }
        onDeselect={(item) =>
          props.setSelected((curr) =>
            curr.filter((filter) => filter !== item.value)
          )
        }>
        <TextDropdownButton>
          {props.columnTitle.toUpperCase()}
        </TextDropdownButton>
      </SelectMenu>
    </Table.TextHeaderCell>
  );
};

// Closure over ServerMenuItem so we can pass in a list of online servers
// in order to set the status indicator when the server is online
export const makeStatusMenuItem = (onlineServers: string[]) => {
  return function (props: any) {
    return (
      <StatusMenuItem
        {...props}
        color={
          onlineServers.includes(props.label) ? "success" : undefined
        }></StatusMenuItem>
    );
  };
};

// Use the default Option menu item but augment it with a status indicator.
// Use "any" because evergreen doesn't export a type for the the <Option/> props.
// See: https://github.com/segmentio/evergreen/blob/674d9a1d31786e7d00088b8553f353d7fee59732/src/select-menu/src/OptionsList.js#L26
export const StatusMenuItem = (props: any) => {
  return (
    <Option {...props}>
      <StatusIndicator color={props.color}></StatusIndicator>
      {props.icon && <Image src={props.icon} width={24} marginRight={8} />}
      {props.label}
    </Option>
  );
};

export default SelectMenuHeader;
