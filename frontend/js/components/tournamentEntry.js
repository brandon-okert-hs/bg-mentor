import React from "react"
import styled from "styled-components"
import missingAvatar from "../image/missingavatar"
import Select from "./basiccomponents/select"

import { COLOR_PRIMARY, COLOR_COMPLIMENT, COLOR_COMPLIMENT_HOVER } from "../colors"

const TournamentEntry = styled.div`
  display: flex;
  flex-flow: row nowrap;
  justify-content: flex-start;
  align-items: center;

  width: 800px;
  height: 120px;

  margin: 10px 0;

  background: ${COLOR_COMPLIMENT};
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

const EntrantColumn = styled.div`
  flex: 1 1;
  width: 100%;
  height: 100%;

  box-sizing: border-box;

  display: flex;
  flex-flow: column nowrap;
  justify-content: space-between;
  align-items: center;
`

const EntrantRow = styled.div`
  display: flex;
  flex-flow: row nowrap;
  justify-content: flex-start;
  align-items: center;

  height: 100%;
  width: 100%;

  margin: 0;
`

const BansArea = styled.div`
  display: flex;
  flex-flow: row wrap;
  justify-content: flex-start;
  align-items: center;

  height: 100%;
  width: 100%;
`

const LockIcon = styled.img`
  flex: 0 0 auto;
  height: 50px;
  width: 50px;
  box-sizing: border-box;

  cursor: ${p => p.isDisabled ? "not-allowed" : "none"};
  opacity: 1;

  &:hover {
    cursor: ${p => p.isDisabled ? "not-allowed" : "pointer"};
    opacity: 0.6;
  }

  &:disabled {
    cursor: ${p => p.isDisabled ? "not-allowed" : "none"};
  }
`

const ReadyIcon = styled.img`
  flex: 0 0 auto;
  height: 30px;
  width: 30px;
  padding: 5px 5px 5px 0;

  cursor: ${p => p.isDisabled ? "not-allowed" : "none"};
  opacity: 1;

  &:hover {
    cursor: ${p => p.isDisabled ? "not-allowed" : "pointer"};
    opacity: 0.6;
  }

  &:disabled {
    cursor: ${p => p.isDisabled ? "not-allowed" : "none"};
  }
`

const ReadyArea = ({ isDisabled, isReady, onClick }) =>
  <ReadyIcon isDisabled={isDisabled} onClick={isDisabled ? null : onClick} src={isReady ? "/static/ready_icon.png" : "/static/not_ready_icon.png"} />

const LockArea = ({ isDisabled, isLocked, onClick }) =>
  <ReadyIcon isDisabled={isDisabled} onClick={isDisabled ? null : onClick} src={isLocked ? "/static/lock_icon.png" : "/static/unlock_icon.png"} />

const Avatar = styled.img`
  flex: 0 0 auto;
  height: 40px;
  width: 40px;

  border-radius: 24px;
`

const put = (id, startLoading, body) => {
  if (body.config.member1 && typeof body.config.member1 === "object") {
    body.config.member1 = body.config.member1.id
  }
  if (body.config.member2 && typeof body.config.member2 === "object") {
    body.config.member2 = body.config.member2.id
  }

  startLoading()
  return fetch(`/api/dabEntry/${id}`, {
    credentials: "include",
    method: "PUT",
    body: JSON.stringify(body)
  }).then(r => r.status === 200 ? r.json() : undefined)
}

const raceItems = [{ label: "zerg", value: "zerg" }, { label: "terran", value: "terran" }, { label: "protoss", value: "protoss" }]
const banNumItems = ["None", 1, 2, 3, 4, 5, 6].map(i => ({ label: i === "None" ? "No bans" : `${i} bans`, value: i }))

export default ({ units = [], members, loggedInMember, creator, startLoading, updateEntry, id, isLocked, tournamentId, config = {}}) => {

  const unitItems = {
    terran: units.filter(u => u.dabCanBan && u.race === "terran").map(u => ({ label: u.name, value: u.name })),
    zerg: units.filter(u => u.dabCanBan && u.race === "zerg").map(u => ({ label: u.name, value: u.name })),
    protoss: units.filter(u => u.dabCanBan && u.race === "protoss").map(u => ({ label: u.name, value: u.name })),
  }

  const memberItems = [{ label: "No Player", value: "None" }].concat(members.map(m => ({ label: m.name, value: m.id })))

  const {
    id: configId,
    member1, member1Race, member1NumBans, member1Confirmed, member1Ban1, member1Ban2, member1Ban3, member1Ban4, member1Ban5, member1Ban6,
    member2, member2Race, member2NumBans, member2Confirmed, member2Ban1, member2Ban2, member2Ban3, member2Ban4, member2Ban5, member2Ban6,
  } = config

  const isCreator = loggedInMember.id === creator.id
  const isMember1 = member1 ? loggedInMember.id === member1.id : false
  const isMember2 = member2 ? loggedInMember.id === member2.id : false
  const isBothConfirmed = member1Confirmed && member2Confirmed

  const onClickLock = () => put(id, startLoading, { config, isLocked: !isLocked }).then(entry => {
    updateEntry()
  })

  const onSelectMember1 = m => {
    const memberId = parseInt(m.target.value, 10)
    put(id, startLoading, { isLocked, config: { ...config, member1: memberId === "None" ? null : memberId }}).then(updateEntry)
  }

  const onSelectMember2 = m => {
    const memberId = parseInt(m.target.value, 10)
    put(id, startLoading, { isLocked, config: { ...config, member2: memberId === "None" ? null : memberId } }).then(updateEntry)
  }

  const onSelectRace1 = m => {
    put(id, startLoading, { isLocked, config: { ...config, member1Race: m.target.value } }).then(updateEntry)
  }

  const onSelectRace2 = m => {
    put(id, startLoading, { isLocked, config: { ...config, member2Race: m.target.value } }).then(updateEntry)
  }

  const onSelectNumBans1 = m => {
    put(id, startLoading, { isLocked, config: { ...config, member1NumBans: parseInt(m.target.value === "None" ? 0 : m.target.value, 10) } }).then(updateEntry)
  }

  const onSelectNumBans2 = m => {
    put(id, startLoading, { isLocked, config: { ...config, member2NumBans: parseInt(m.target.value === "None" ? 0 : m.target.value, 10) } }).then(updateEntry)
  }

  const createOnSelectBan = banTarget => m =>
    put(id, startLoading, { isLocked, config: { ...config, [banTarget]: m.target.value === "None" ? null : m.target.value } }).then(updateEntry)
  
  const member1BansSelectors = [
    <Select hiddenValue={member1Ban1 ? "<chosen>" : "<unchosen>"} isHidden={!isCreator && !isMember1 && !isBothConfirmed} isDisabled={!isCreator && (!isMember1 || isLocked)} key={1} onClick={createOnSelectBan("member1Ban1")} items={unitItems[member2Race || "zerg"]} value={member1Ban1 || "None"} />,
    <Select hiddenValue={member1Ban2 ? "<chosen>" : "<unchosen>"} isHidden={!isCreator && !isMember1 && !isBothConfirmed} isDisabled={!isCreator && (!isMember1 || isLocked)} key={2} onClick={createOnSelectBan("member1Ban2")} items={unitItems[member2Race || "zerg"]} value={member1Ban2 || "None"} />,
    <Select hiddenValue={member1Ban3 ? "<chosen>" : "<unchosen>"} isHidden={!isCreator && !isMember1 && !isBothConfirmed} isDisabled={!isCreator && (!isMember1 || isLocked)} key={3} onClick={createOnSelectBan("member1Ban3")} items={unitItems[member2Race || "zerg"]} value={member1Ban3 || "None"} />,
    <Select hiddenValue={member1Ban4 ? "<chosen>" : "<unchosen>"} isHidden={!isCreator && !isMember1 && !isBothConfirmed} isDisabled={!isCreator && (!isMember1 || isLocked)} key={4} onClick={createOnSelectBan("member1Ban4")} items={unitItems[member2Race || "zerg"]} value={member1Ban4 || "None"} />,
    <Select hiddenValue={member1Ban5 ? "<chosen>" : "<unchosen>"} isHidden={!isCreator && !isMember1 && !isBothConfirmed} isDisabled={!isCreator && (!isMember1 || isLocked)} key={5} onClick={createOnSelectBan("member1Ban5")} items={unitItems[member2Race || "zerg"]} value={member1Ban5 || "None"} />,
    <Select hiddenValue={member1Ban6 ? "<chosen>" : "<unchosen>"} isHidden={!isCreator && !isMember1 && !isBothConfirmed} isDisabled={!isCreator && (!isMember1 || isLocked)} key={6} onClick={createOnSelectBan("member1Ban6")} items={unitItems[member2Race || "zerg"]} value={member1Ban6 || "None"} />,
  ]

  const member2BansSelectors = [
    <Select hiddenValue={member2Ban1 ? "<chosen>" : "<unchosen>"} isHidden={!isCreator && !isMember2 && !isBothConfirmed} isDisabled={!isCreator && (!isMember2 || isLocked)} key={1} onClick={createOnSelectBan("member2Ban1")} items={unitItems[member1Race || "zerg"]} value={member2Ban1 || "None"} />,
    <Select hiddenValue={member2Ban1 ? "<chosen>" : "<unchosen>"} isHidden={!isCreator && !isMember2 && !isBothConfirmed} isDisabled={!isCreator && (!isMember2 || isLocked)} key={2} onClick={createOnSelectBan("member2Ban2")} items={unitItems[member1Race || "zerg"]} value={member2Ban2 || "None"} />,
    <Select hiddenValue={member2Ban1 ? "<chosen>" : "<unchosen>"} isHidden={!isCreator && !isMember2 && !isBothConfirmed} isDisabled={!isCreator && (!isMember2 || isLocked)} key={3} onClick={createOnSelectBan("member2Ban3")} items={unitItems[member1Race || "zerg"]} value={member2Ban3 || "None"} />,
    <Select hiddenValue={member2Ban1 ? "<chosen>" : "<unchosen>"} isHidden={!isCreator && !isMember2 && !isBothConfirmed} isDisabled={!isCreator && (!isMember2 || isLocked)} key={4} onClick={createOnSelectBan("member2Ban4")} items={unitItems[member1Race || "zerg"]} value={member2Ban4 || "None"} />,
    <Select hiddenValue={member2Ban1 ? "<chosen>" : "<unchosen>"} isHidden={!isCreator && !isMember2 && !isBothConfirmed} isDisabled={!isCreator && (!isMember2 || isLocked)} key={5} onClick={createOnSelectBan("member2Ban5")} items={unitItems[member1Race || "zerg"]} value={member2Ban5 || "None"} />,
    <Select hiddenValue={member2Ban1 ? "<chosen>" : "<unchosen>"} isHidden={!isCreator && !isMember2 && !isBothConfirmed} isDisabled={!isCreator && (!isMember2 || isLocked)} key={6} onClick={createOnSelectBan("member2Ban6")} items={unitItems[member1Race || "zerg"]} value={member2Ban6 || "None"} />,
  ]

  const onClickReady1 = () =>
    put(id, startLoading, { isLocked, config: { ...config, member1Confirmed: !member1Confirmed } }).then(updateEntry)

  const onClickReady2 = () =>
    put(id, startLoading, { isLocked, config: { ...config, member2Confirmed: !member2Confirmed } }).then(updateEntry)

  return (
    <TournamentEntry>
      <LockArea isDisabled={!isCreator} isLocked={isLocked} onClick={onClickLock} />
      <EntrantColumn>
        <EntrantRow>
          <Avatar src={member1 && member1.avatarUrl || missingAvatar} />
          <Select isDisabled={!isCreator} onClick={onSelectMember1} items={memberItems} value={member1 ? member1.id : null} />
          <Select isDisabled={!isCreator && (!isMember1 || isLocked)} onClick={onSelectRace1} items={raceItems} value={member1Race || "zerg"} />
          <Select isDisabled={!isCreator} onClick={onSelectNumBans1} items={banNumItems} value={member1NumBans || "None"} />
          <BansArea>
            {member1BansSelectors.slice(0, member1NumBans || 0)}
          </BansArea>
          <ReadyArea isDisabled={!isCreator && (!isMember1 || isLocked)} isReady={member1Confirmed} onClick={onClickReady1} />
        </EntrantRow>
        <EntrantRow>
          <Avatar src={member2 && member2.avatarUrl || missingAvatar} />
          <Select isDisabled={!isCreator} onClick={onSelectMember2} items={memberItems} value={member2 ? member2.id : null} />
          <Select isDisabled={!isCreator && (!isMember2 || isLocked)} onClick={onSelectRace2} items={raceItems} value={member2Race || "zerg"} />
          <Select isDisabled={!isCreator} onClick={onSelectNumBans2} items={banNumItems} value={member2NumBans || "None"} />
          <BansArea>
            {member2BansSelectors.slice(0, member2NumBans || 0)}
          </BansArea>
          <ReadyArea isDisabled={!isCreator && (!isMember2 || isLocked)} isReady={member2Confirmed} onClick={onClickReady2} />
        </EntrantRow>
      </EntrantColumn>
    </TournamentEntry>
  )
}
