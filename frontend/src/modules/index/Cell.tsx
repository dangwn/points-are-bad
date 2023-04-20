import React from 'react';

interface CellProps {
  topText: string|number,
  bottomText: string,
  topSize: string,
  bottomSize: string
  className: any
}

const Cell: React.FC<CellProps> = ({topText, bottomText, topSize, bottomSize, className}) => {
  return (
    <div className={className}>
      <div style={{'fontWeight':'bold','fontSize':topSize}}>{topText}</div>
      <div style={{'fontSize':bottomSize}}>{bottomText}</div>
    </div>
  )
}

export default Cell;