import React from "react"
import ReactDOM from "react-dom"
import styled from "styled-components"

import Button, { BUTTON_TYPE_SECONDARY, BUTTON_TYPE_PRIMARY } from "./basiccomponents/button"
import { COLOR_UI } from './colors'

const Navbar = styled.div`
  display: flex;
  flex-flow: row nowrap;
  justify-content: flex-end;
  align-items: center;
  min-height: 40px;

  background: ${COLOR_UI};
`

const NavIcon = styled.img`
  height: 34px;
  margin-right: 4px;
`

export default ({ rightItems = [], leftItems = [] }) => {
  const leftItemList = leftItems.map(item =>
    <Button
      key={item.name}
      onClick={() => open(item.url, "_self")}
      type={BUTTON_TYPE_PRIMARY}
      height="40px"
    >
      {item.image ? <NavIcon src={item.image} /> : null}
      {item.name}
    </Button> 
  )

  const rightItemList = rightItems.map((item, i) =>
    <Button
      key={item.name}
      style={i === 0 ? {"marginLeft": "auto"} : {}}
      onClick={() => open(item.url, "_self")}
      type={BUTTON_TYPE_PRIMARY}
      height="40px"
      >
        {item.image ? <NavIcon src={item.image} /> : null}
        {item.name}
    </Button>
  )

  return <Navbar>{leftItemList}{rightItemList}</Navbar>
}
