import { Dispatch } from "@reduxjs/toolkit";
import clsx from "clsx";
import {
  Badge,
  Menu,
  Pane,
  Popover,
  Position,
  Spinner,
  Text
} from "evergreen-ui";
import { Eye, EyeSlash, MagnifyingGlass, Trash } from "phosphor-react";
import React from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  deleteHistoryItem,
  HistoryItem,
  selectHistory
} from "../../state/historySlice";
import { setActiveItem } from "../../state/stateSlice";
import { RootState } from "../../state/store";

const SearchHistory: React.FC = () => {
  const history = useSelector(selectHistory);
  const activeTS =
    useSelector((store: RootState) => store.state.activeItem?.timestamp) ?? -1;
  const dispatch = useDispatch();

  return (
    <>
      {history.length > 0 ? (
        history.map((item: HistoryItem) => (
          <HistoryCard
            activeTS={activeTS}
            key={item.timestamp.toString()}
            item={item}
            dispatch={dispatch}
          />
        ))
      ) : (
        <p className="text-center my-4 text-gray-500 text-sm">
          History is a mystery.
        </p>
      )}
    </>
  );
};

type Props = {
  activeTS: number;
  item: HistoryItem;
  dispatch: Dispatch<any>;
};

const HistoryCard: React.FC<Props> = ({ activeTS, item, dispatch }: Props) => {
  const isActive = activeTS === item.timestamp;

  return (
    <Popover
      position={Position.BOTTOM}
      content={
        <Menu>
          {!isActive && (
            <Menu.Item
              icon={<Eye color="black" size={18} weight="bold" />}
              onClick={() => dispatch(setActiveItem(item))}>
              Show
            </Menu.Item>
          )}
          {isActive && (
            <Menu.Item
              icon={<EyeSlash color="black" size={18} weight="bold" />}
              onClick={() => dispatch(setActiveItem(null))}>
              Hide
            </Menu.Item>
          )}
          <Menu.Item
            icon={<Trash color="#d45050" size={18} weight="bold" />}
            onClick={() => dispatch(deleteHistoryItem(item.timestamp))}
            intent="danger">
            Delete
          </Menu.Item>
        </Menu>
      }>
      <div
        className={clsx(
          "border cursor-pointer p-2 m-1 rounded hover:bg-hover-blue flex justify-between items-center",
          { "shadow border-active-text-blue": isActive }
        )}>
        <MagnifyingGlass size={18} weight="bold" />
        <Text
          width={140}
          marginLeft={24}
          flex={1}
          whiteSpace="nowrap"
          display="block"
          fontWeight={500}
          overflow="hidden"
          textOverflow="ellipsis">
          {item.query}
        </Text>
        {!item.results?.length && !item.errors?.length ? (
          <Spinner size={20} marginRight={5} />
        ) : (
          <Badge color="blue">{`${item.results?.length} RESULTS`}</Badge>
        )}
      </div>
    </Popover>
  );
};

export default SearchHistory;
