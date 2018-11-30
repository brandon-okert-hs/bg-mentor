import React from "react"
import ReactDOM from "react-dom"
import styled from "styled-components"

import { COLOR_PRIMARY } from "./colors"

const UserWelcome = styled.div`
  background: ${COLOR_PRIMARY};

  height: 40px;
  padding: 2px 20px;
  box-sizing: border-box;
  margin: 36px 0;

  display: flex;
  flex-flow: row nowrap;
  justify-content: flex-start;
  align-items: center;
`

const Name = styled.div`
  display: flex;
  flex-flow: column nowrap;
  justify-content: center;
  height: 40px;
  padding-right: 10px;

  color: white;

  word-wrap: break-word;
  overflow: hidden;
`

const Avatar = styled.img`
  height: 50px;
  width: 50px;

  border-radius: 24px;
`

export default ({ avatarUrl, name }) => (
  <UserWelcome>
    <Name>{name}</Name>
    <Avatar src={avatarUrl} />
  </UserWelcome>
)


