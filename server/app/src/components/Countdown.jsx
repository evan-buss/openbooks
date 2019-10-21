import React from 'react';
import './Countdown.css'

export default function Countdown(props) {
    return (
        <div id="countdown">
            <h1>Please wait</h1>
            <h2>{props.time} seconds</h2>
            <h1>before searching.</h1>
        </div>
    )
}