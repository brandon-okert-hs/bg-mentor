import React from "react"
import ReactDOM from "react-dom"
import styled from "styled-components"
import moment from "moment"

import { COLOR_PRIMARY, COLOR_COMPLIMENT, COLOR_COMPLIMENT_HOVER } from "../colors"

const Tournament = styled.div`
  display: flex;
  flex-flow: row nowrap;
  justify-content: flex-start;
  align-items: center;

  width: 800px;
  height: 95px;

  margin: 10px 0;

  background: ${COLOR_COMPLIMENT};
  cursor: none;

  &:hover {
    background: ${COLOR_COMPLIMENT_HOVER};
    cursor: pointer;
  }
`

const DetailArea = styled.div`
  width: 100%;
  height: 40px;

  display: flex;
  flex-flow: row nowrap;
  justify-content: flex-start;
  align-items: center;
`

const DetailIcon = styled.img`
  flex: 0 0 auto;
  width: 30px;
  height: 30px;
`

const DetailValue = styled.span`
  flex: 0 0 auto;
  font-size: 16px;
  color: white;

  padding: 5px;
  
  width: 60%;
`

const Detail = ({ iconUrl, children }) => (
  <DetailArea>
    <DetailIcon src={iconUrl} />
    <DetailValue>{children}</DetailValue>
  </DetailArea>
)

const DetailBox = styled.div`
  display: flex;
  flex-flow: column nowrap;
  justify-content: space-between;
  align-items: center;

  padding: 5px 5px 5px 0;
  box-sizing: border-box;

  width: 100%;
  height: 100%;
`

const DetailBoxContainer = styled.div`
  flex: 1 1 auto;
  width: 100%;
  height: 100%;

  box-sizing: border-box;

  display: flex;
  flex-flow: row nowrap;
  justify-content: space-between;
  align-items: center;
`

const Logo = styled.img`
  flex: 0 0 auto;
  height: 95px;
  width: 95px;
  padding: 5px;
  box-sizing: border-box;
`

export default ({ onClick, creator, name, logoLink, startDate, checkinDate, type = "custom", entries = [] }) => {

  let start = ""
  let checkin = ""
  if (startDate) {
    const hoursUntil = startDate.diff(moment.utc(), "hours")
    if (hoursUntil > 48) {
      start = `In about ${Math.round(hoursUntil / 24)} days`
    } else if (hoursUntil > 1) {
      start = `In about ${hoursUntil} hours`
    } else {
      start = `In ${hoursUntil * 60} minutes`
    }

    if (checkinDate) {
      const checkinMinutes = startDate.diff(checkinDate, "minutes")
      checkin = `${checkinMinutes} minutes before`
    }
  }

  return (
    <Tournament onClick={onClick} >
      <Logo src={logoLink || "/static/borngosu_logo.png"} />
      <DetailBoxContainer>
        <DetailBox>
          <Detail iconUrl={"/static/tournament_icon.png"}>{name}</Detail>
          <Detail iconUrl={"/static/manager_icon.png"}>{creator.name}</Detail>
        </DetailBox>
        <DetailBox>
          <Detail iconUrl={"/static/player_icon.png"}>{entries.length}</Detail>
          <Detail iconUrl={"/static/rules_icon.png"}>{type.toUpperCase()}</Detail>
        </DetailBox>
        <DetailBox>
          <Detail iconUrl={"/static/clock_icon.png"}>{start}</Detail>
          <Detail iconUrl={"/static/alarm_icon.png"}>{checkin}</Detail>
        </DetailBox>
      </DetailBoxContainer>
    </Tournament>
  )
}


