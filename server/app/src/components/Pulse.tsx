import { Position, Tooltip } from 'evergreen-ui';
import React from 'react';
import "./Pulse.css";

type Props = {
    disabled?: boolean;
}

const Pulse: React.FC<Props> = ({ disabled = false }) => {
    return (
        <Tooltip position={Position.BOTTOM} content={`OpenBooks server ${disabled ? 'disconnected' : 'connected'}.`}>
            <div className="center">
                <div className={`circle pulse ${disabled && 'disabled'}`}></div>
            </div>
        </Tooltip>
    );
}

export default Pulse;
