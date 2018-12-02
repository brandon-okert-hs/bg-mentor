import React from "react"
import ReactDOM from "react-dom"
import styled from "styled-components"

import {
  COLOR_UI,
  COLOR_UI_HOVER,
} from "../../colors"

const Select = styled.select`
  max-width: 120px;

  display: flex;
  flex-flow: row nowrap;
  justify-content: center;
  align-items: center;

  border: none;
  border-radius: none;

  padding: 6px;
  margin: 4px;
  box-sizing: border-box;

  font-size: 16px;
  color: white;
  cursor: none;

  background: ${COLOR_UI};

  &:hover {
    cursor: pointer;
    background: ${COLOR_UI_HOVER};
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
`

export default ({ items, value, onClick, isDisabled = false, isHidden = false, hiddenValue = "<hidden>" }) => isHidden ? (
  <Select disabled={isDisabled} value={"hidden"}><option value={"hidden"}>{hiddenValue}</option></Select>
) : (
  <Select disabled={isDisabled} value={value || -1} onChange={onClick}>
    {items.map(i => <option key={i.value} value={i.value}>{i.label}</option>)}
  </Select>
)
