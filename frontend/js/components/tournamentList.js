import React from "react"
import ReactDOM from "react-dom"
import styled from "styled-components"
import missingAvatar from "../image/missingavatar"

import Loader from './basiccomponents/loader'
import Tournament from "./tournament"

const TournamentList = styled.div`
  height: 100%;
  width: 100%;

  display: flex;
  flex-flow: column nowrap;
  justify-content: flex-start;
  align-items: center;
`

export default ({ tournaments }) => (
  <TournamentList>
    {tournaments.map((t, i) => <Tournament key={i} {...t} />)}
    {tournaments.length === 0 ? <Loader /> : undefined}
  </TournamentList>
)
