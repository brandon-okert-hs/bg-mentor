import React from "react"
import ReactDOM from "react-dom"
import styled from "styled-components"

import Button, { BUTTON_TYPE_SECONDARY, BUTTON_TYPE_PRIMARY } from "./basiccomponents/button"
import { COLOR_UI, COLOR_COMPLIMENT } from '../colors'

const leftNavItems = [{
  name: "Born Gosu Gaming",
  url: "/",
  image: "/static/borngosu_logo.png",
}, {
  name: "Tournaments",
  url: "/tournaments",
  image: "/static/challonge_logo.png",
}]
const rightNavItems = [{
  name: "Login",
  url: "/auth/login",
  image: "",
}, {
  name: "Logout",
  url: "/auth/logout",
  image: "",
}]

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

export default () => {
  const leftItemList = leftNavItems.map(item =>
    <Button
      key={item.name}
      onClick={() => open(item.url, "_self")}
      type={BUTTON_TYPE_PRIMARY}
      selected={window.location.pathname === item.url}
      height="40px"
    >
      {item.image ? <NavIcon src={item.image} /> : null}
      {item.name}
    </Button> 
  )

  const rightItemList = rightNavItems.map((item, i) =>
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
