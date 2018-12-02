import React from "react"
import ReactDOM from "react-dom"
import styled from "styled-components"
import missingAvatar from "../image/missingavatar"

import Loader from './basiccomponents/loader'
import TournamentEntry from "./tournamentEntry"

const TournamentEntryList = styled.div`
  height: 100%;
  width: 100%;

  display: flex;
  flex-flow: column nowrap;
  justify-content: flex-start;
  align-items: center;
`

const AddEntryButton = styled.img`
  width: 60px;
  height: 60px;

  cursor: none;

  &:hover {
    cursor: pointer;
    opacity: 0.7;
  }
`

export default ({ units, tournament, members, loggedInMember, updateEntry, startLoading, isLoading }) => {
  
  const onClickAddEntry = () => {
    startLoading()
    return fetch("/api/dabEntry/", {
      credentials: "include",
      method: "POST",
      body: JSON.stringify({
        tournamentID: tournament.id,
      })
    }).then(r => r.status === 200 ? r.json() : undefined).then(updateEntry)
  }

  const isCreator = loggedInMember && tournament && tournament.creator && loggedInMember.id === tournament.creator.id

  return (
    <TournamentEntryList>
      {tournament ? tournament.entries.map(e => <TournamentEntry loggedInMember={loggedInMember} creator={tournament.creator} units={units} members={members} startLoading={startLoading} updateEntry={updateEntry} key={e.id} {...e} />) : undefined}
      {!tournament || isLoading ? <Loader /> : undefined}
      {isCreator ? <AddEntryButton src="/static/add_entry_icon.png" onClick={onClickAddEntry} /> : null}
    </TournamentEntryList>
  )
}
