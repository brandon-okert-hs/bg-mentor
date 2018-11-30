import React from "react"
import ReactDOM from "react-dom"
import styled from "styled-components"

import {
  COLOR_PRIMARY,
  COLOR_COMPLIMENT,
  COLOR_COMPLIMENT_HOVER,
  COLOR_OUTLINE,
  COLOR_UI,
  COLOR_UI_HOVER,
  COLOR_UI_ACTIVE,
} from "../colors"

export const BUTTON_TYPE_PRIMARY = "BUTTON_TYPE_PRIMARY"
export const BUTTON_TYPE_SECONDARY = "BUTTON_TYPE_SECONDARY"

const typeToDefaultColor = type => {
  switch (type) {
    case BUTTON_TYPE_PRIMARY: return COLOR_UI
    case BUTTON_TYPE_SECONDARY: return COLOR_COMPLIMENT
  }
}

const typeToHoverColor = type => {
  switch (type) {
    case BUTTON_TYPE_PRIMARY: return COLOR_UI_HOVER
    case BUTTON_TYPE_SECONDARY: return COLOR_COMPLIMENT_HOVER
  }
}

export default styled.div`
  width: ${p => p.width};
  height: ${p => p.height};

  display: flex;
  flex-flow: row nowrap;
  justify-content: center;
  align-items: center;

  padding: 6px;
  box-sizing: border-box;

  font-size: 16px;
  color: white;
  cursor: none;

  background: ${p => typeToDefaultColor(p.type)};

  &:hover {
    cursor: pointer;
    background: ${p => typeToHoverColor(p.type)};
  }
`
