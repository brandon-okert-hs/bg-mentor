import React from "react"
import ReactDOM from "react-dom"
import styled from "styled-components"

import NavBar from './components/navbar'
import TournamentList from './components/tournamentList'

import moment from "moment"

document.body.style.margin = 0;

const App = styled.div`
  width: 100%;
  height: 100%;
  font-family: 'Open Sans', sans-serif;
  font-size: 14px;

  display: flex;
  flex-flow: column nowrap;
  justify-content: flex-start;
  align-items: stretch;
`

const app = document.createElement('div')
app.style.width = "100vw"
app.style.height = "100vh"
document.body.appendChild(app)

const leftNavItems = [{
  name: "Born Gosu Gaming",
  url: "/",
  image: "/static/borngosu_logo.png",
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

class AppContainer extends React.Component {
  constructor() {
    super()
    this.state = {
      loggedInMember: null,
      tournaments: undefined,
    }

    fetch("/api/member/me", { credentials: "include" }).then(r => {
      if (r.status === 200) {
        return r.json()
      }
    }).then(member => {
      if (member) {
        this.setState({
          loggedInMember: member,
        })
      }
    })

    fetch("/api/tournament", { credentials: "include" }).then(r => {
      if (r.status === 200) {
        return r.json()
      }
    }).then(tournaments => {
      if (tournaments) {
        tournaments.forEach(t => {
          if (t.startDate) {
            t.startDate = moment.utc(t.startDate, "YYYY-MM-DD HH:mm:ss")
          }
          if (t.checkinDate) {
            t.checkinDate = moment.utc(t.checkinDate, "YYYY-MM-DD HH:mm:ss")
          }
        })

        const entryPromises = tournaments.map(t =>
          fetch(`/api/dabEntry?tournamentID=${t.id}`, { credentials: "include" })
            .then(r => r.status === 200 ? r.json() : undefined)
            .then(entries => ({ tid: t.id, entries: entries || [] }))
        )
        
        const memberPromise = fetch(`/api/member`, { credentials: "include" })
          .then(r => r.status === 200 ? r.json() : undefined)

        Promise.all(entryPromises).then(entryLists => {
          tournaments.forEach(t => t.entries = entryLists.find(e => e.tid === t.id).entries)

          memberPromise.then(members => {
            tournaments.forEach(t => {
              t.creator = members.find(m => m.id === t.creator)
              t.entries.forEach(e => {
                if (e.config.member1) e.config.member1 = members.find(m => m.id === e.config.member1)
                if (e.config.member2) e.config.member2 = members.find(m => m.id === e.config.member2)
              })
            })
            this.setState({ tournaments })
          })
        })
      }
    })
  }

  render() {
    return (
      <App>
        <NavBar />
        <TournamentList tournaments={this.state.tournaments || []} isLoading={!Array.isArray(this.state.tournaments)} />
      </App>
    )
  }
}

ReactDOM.render(<AppContainer />, app)
