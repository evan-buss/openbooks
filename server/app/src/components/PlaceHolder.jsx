import React from 'react';
import Title from 'antd/es/typography/Title';
import image from '../book_person.svg'

const svgStyle = {
  width: "40%"
}

const placeholderStyle = {
  marginTop: 120,
  display: "flex",
  flexFlow: "column nowrap",
  justifyContent: "center",
  alignItems: "center"
}

export default function PlaceHolder() {
  return (
    <div style={placeholderStyle}>
      <img style={svgStyle} src={image} alt="placeholder" />
      <Title level={2}>Search a book to get started</Title>
    </div>
  );
}