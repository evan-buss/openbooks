import clsx from "clsx";
import { Tooltip, Position } from "evergreen-ui";
import { Bell } from "phosphor-react";
import React from "react";

interface PulseProps {
  enabled: boolean;
  onClick: React.MouseEventHandler<HTMLDivElement>;
}

const Pulse = ({ enabled, onClick }: PulseProps) => {
  return (
    <Tooltip
      showDelay={500}
      position={Position.BOTTOM}
      content={`OpenBooks server ${enabled ? "connected" : "disconnected"}.`}>
      <div onClick={onClick} className="flex justify-evenly items-center">
        <Bell
          size={18}
          weight="bold"
          className={clsx([
            "rounded-full cursor-pointer hover:text-blue-800",
            {
              "text-active-text-blue animate-custom-pulse hover:text-gray-500":
                enabled
            },
            { "animate-none text-gray-500": !enabled }
          ])}
        />
      </div>
    </Tooltip>
  );
};

export default Pulse;
