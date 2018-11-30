import React from "react"
import ReactDOM from "react-dom"
import styled from "styled-components"

import NavBar from './navbar'
import MainMenu from './mainmenu'

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
      loggedInMember: null
    }

    fetch("member/me", { credentials: "include" }).then(r => {
      if (r.status === 200) {
        return r.json()
      }
    }).then(member => {
      if (member) {
        this.setState({
          loggedInMember: member,
        });
      }
    }).catch(alert);
  }
  
  render() {
    return (
      <App>
        <NavBar leftItems={leftNavItems} rightItems={rightNavItems} />
        <MainMenu loggedInMember={this.state.loggedInMember} />
      </App>
    )
  }
}

ReactDOM.render(<AppContainer />, app)
