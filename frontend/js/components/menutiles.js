import React from "react"
import ReactDOM from "react-dom"
import styled from "styled-components"

import { COLOR_PRIMARY, COLOR_COMPLIMENT, COLOR_COMPLIMENT_HOVER } from "../colors"

const MenuTiles = styled.div`
  display: flex;
  flex-flow: row nowrap;
  justify-content: space-around;
  align-items: center;

  width: 100%;
  max-width: 800px;
`

const MenuTileArea = styled.div`
  background: ${COLOR_COMPLIMENT};

  height: 200px;
  width: 120px;

  display: flex;
  flex-flow: column nowrap;
  justify-content: flex-start;
  align-items: stretch;

  cursor: none;

  &:hover {
    background: ${COLOR_COMPLIMENT_HOVER};
    cursor: pointer;
  }
`

const Name = styled.div`
  display: flex;
  flex-flow: row nowrap;
  justify-content: center;
  align-items: center;
  height: 40px;

  color: white;

  word-wrap: break-word;
  overflow: hidden;
`

const CallToActionArea = styled.div`
  height: 120px;

  display: flex;
  flex-flow: row nowrap;
  justify-content: center;
  align-items: center;
`

const CTAImage = styled.img`
  width: 80px;
`

const CallToAction= ({ imageUrl }) => (
  <CallToActionArea>
    {imageUrl ? <CTAImage src={imageUrl}/> : null }
  </CallToActionArea>
)

const Disclaimer = styled.div`
  background: ${COLOR_PRIMARY};
  display: flex;
  flex-flow: row nowrap;
  justify-content: center;
  align-items: center;
  height: 40px;

  color: white;

  word-wrap: break-word;
  overflow: hidden;
`

const MenuTile = ({ name, link, image, isExternal, disclaimer, className }) => (
  <MenuTileArea className={className} onClick={() => link ? open(link, isExternal ? undefined : "_self") : undefined}>
    <Name>{name}</Name>
    <CallToAction imageUrl={image} />
    {disclaimer ? <Disclaimer>{disclaimer}</Disclaimer> : null}
  </MenuTileArea>
)

export default ({ tiles }) => (
  <MenuTiles>
    {tiles.map(t => <MenuTile key={t.name} {...t} />)}
  </MenuTiles>
)


