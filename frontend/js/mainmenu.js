import React from "react"
import ReactDOM from "react-dom"
import styled from "styled-components"
import missingAvatar from "./image/missingavatar"

import UserWelcome from "./userwelcome"
import MenuTiles from "./menutiles"

const MainMenu = styled.div`
  background: "#efefef";
  height: 100%;

  display: flex;
  flex-flow: column nowrap;
  justify-content: flex-start;
  align-items: center;
`

export default ({ loggedInMember }) => {
  const name = loggedInMember ? loggedInMember.name : "guest"
  const avatarUrl = loggedInMember ? loggedInMember.avatarUrl : missingAvatar

  const tiles = [{
    name: "Discord",
    link: undefined,
    image: "/static/discord_logo.png",
    disclaimer: "Coming Soon!"
  }, {
    name: "Liquipedia",
    image: "/static/liquipedia_logo.png",
    link: undefined,
    disclaimer: "Coming Soon!"
  }, {
    name: "Tournaments",
    image: "/static/challonge_logo.png",
    link: undefined,
    disclaimer: loggedInMember ? "Coming Soon!" : "Members Only"
  }, {
    name: "Mentor Program",
    image: "/static/grandmaster_logo.png",
    link: undefined,
    disclaimer: loggedInMember ? "Coming Soon!" : "Members Only"
  }, {
    name: "Twitch",
    image: "/static/twitch_logo.png",
    link: undefined,
    disclaimer: "Coming Soon!"
  }]

  return (
    <MainMenu>
      <UserWelcome name={name} avatarUrl={avatarUrl} />
      <MenuTiles tiles={tiles} />
    </MainMenu>
  )
}
