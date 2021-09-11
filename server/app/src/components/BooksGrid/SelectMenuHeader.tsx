import {
  SelectMenu,
  SelectMenuItem,
  Table,
  TextDropdownButton
} from "evergreen-ui";
import React from "react";

interface Props {
  columnTitle: string;
  menuTitle: string;
  options: string[];
  selected: string[];
  setSelected: React.Dispatch<React.SetStateAction<string[]>>;
}

const SelectMenuHeader = (props: Props) => {
  const menuOptions: SelectMenuItem[] = props.options.map((opt) => ({
    label: opt,
    value: opt
  }));

  return (
    <Table.TextHeaderCell flexBasis={100} flexGrow={0} flexShrink={0}>
      <SelectMenu
        isMultiSelect
        title={props.menuTitle}
        options={menuOptions}
        selected={props.selected}
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

export default SelectMenuHeader;
