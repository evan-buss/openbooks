import { Button } from "evergreen-ui";
import { AnimatePresence } from "framer-motion";
import React, { useState } from "react";
import { useSelector } from "react-redux";
import { sendDownload } from "../state/stateSlice";
import { useAppDispatch, RootState } from "../state/store";
import ThreeDotWave from "./BooksGrid/ThreeDotWave";

export function DownloadButton({ book }: { book: string }) {
  const dispatch = useAppDispatch();
  const [disabled, setDisabled] = useState(false);

  const isInFlight = useSelector((state: RootState) =>
    state.state.inFlightDownloads.includes(book)
  );

  // Prevent hitting the same button multiple times
  const onClick = () => {
    if (disabled) return;
    dispatch(sendDownload(book));
    setDisabled(true);
  };

  return (
    <Button
      appearance="primary"
      size="small"
      width="100px"
      disabled={disabled}
      onClick={onClick}>
      <AnimatePresence>
        {isInFlight ? <ThreeDotWave /> : <span>Download</span>}
      </AnimatePresence>
    </Button>
  );
}
