import React from "react"
import ReactDOM from "react-dom"
import styled from "styled-components"

import NavBar from './components/navbar'
import TournamentEntryList from './components/tournamentEntryList'

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
      tournament: undefined,
      isLoading: false,
      members: [],
      units: []
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

    fetch(`/api/unit`, { credentials: "include" }).then(r => r.status === 200 ? r.json() : undefined).then(units => {
      if (units) {
        this.setState({
          units: units,
        })
      }
    })

    const pathElements = window.location.pathname.split("/")
    const tournamentID = pathElements[pathElements.length-1]

    this.updateEntry = () => {
      return fetch(`/api/tournament/${tournamentID}`, { credentials: "include" }).then(r => {
        if (r.status === 200) {
          return r.json()
        }
      }).then(t => {
        if (t) {
          if (t.startDate) {
            t.startDate = moment.utc(t.startDate, "YYYY-MM-DD HH:mm:ss")
          }
          if (t.checkinDate) {
            t.checkinDate = moment.utc(t.checkinDate, "YYYY-MM-DD HH:mm:ss")
          }

          const entryPromise = fetch(`/api/dabEntry?tournamentID=${t.id}`, { credentials: "include" })
            .then(r => r.status === 200 ? r.json() : undefined)
            .then(entries => ({ tid: t.id, entries: entries || [] }))

          const memberPromise = fetch(`/api/member`, { credentials: "include" })
            .then(r => r.status === 200 ? r.json() : undefined)

          entryPromise.then(e => {
            t.entries = e.entries

            memberPromise.then(members => {
              t.creator = members.find(m => m.id === t.creator)
              t.entries.forEach(e => {
                if (e.config.member1) e.config.member1 = members.find(m => m.id === e.config.member1)
                if (e.config.member2) e.config.member2 = members.find(m => m.id === e.config.member2)
              })

              this.setState({ tournament: t, isLoading: false, members: members })
            })
          })
        }
      })
    }

    this.updateEntry()
  }

  render() {
    return (
      <App>
        <NavBar />
        <TournamentEntryList units={this.state.units} startLoading={() => this.setState({ isLoading: true })} updateEntry={this.updateEntry} loggedInMember={this.state.loggedInMember} tournament={this.state.tournament} members={this.state.members} isLoading={this.state.isLoading} />
      </App>
    )
  }
}

ReactDOM.render(<AppContainer />, app)
