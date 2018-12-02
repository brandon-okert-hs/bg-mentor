import React from "react"
import ReactDOM from "react-dom"
import styled from "styled-components"
import { COLOR_COMPLIMENT } from "../../colors"

export default styled.div`
  width: 30px;
  height: 30px;

  margin: 100px;

  background: ${COLOR_COMPLIMENT};

  animation:spin 3s linear infinite;
  @keyframes spin { 100% { -webkit-transform: rotate(360deg); transform:rotate(360deg); } }
`
