import { Pane } from "evergreen-ui";
import styled, { css } from "styled-components";

const TogglePane = styled(Pane)`
    ${props => props.active
        ? css`border-left: 3px solid #1070CA;`
        : css`border-left: 1px solid #E4E7EB;`
    }
`

export default TogglePane;
