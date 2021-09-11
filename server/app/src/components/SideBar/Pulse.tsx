import clsx from "clsx";
import { Tooltip, Position, NotificationsIcon } from "evergreen-ui";
import React from "react";

interface PulseProps {
  enabled: boolean;
  onClick: React.MouseEventHandler<HTMLDivElement>;
}

const Pulse = ({ enabled, onClick }: PulseProps) => {
  return (
    <Tooltip
      position={Position.BOTTOM}
      content={`OpenBooks server ${enabled ? "connected" : "disconnected"}.`}>
      <div onClick={onClick} className="flex justify-evenly items-center">
        <NotificationsIcon
          className={clsx([
            "w-3 h-3 rounded-full cursor-pointer",
            { "text-custom-blue animate-custom-pulse": enabled },
            { "animate-none text-gray-500": !enabled }
          ])}
        />
        {/* <div
          className={clsx([
            "w-2 h-2 rounded-full",
            { "bg-custom-blue animate-custom-pulse": enabled },
            { "animate-none bg-gray-500": !enabled }
          ])}
        /> */}
      </div>
    </Tooltip>
  );
};

export default Pulse;
