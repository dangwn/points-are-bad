export const preventNegativeInputs = (e: React.KeyboardEvent<HTMLInputElement>): void => {
  if (e.key === '-' || e.key === '_') {
    e.preventDefault();
  };
}

export const createPositionString = (position: number|string|null): string => {
  if (position === null){
    return '-'
  };

  const suffixes = ['th', 'st', 'nd', 'rd'];
  const positionNumber = parseInt(position.toString(), 10);
  const suffix = positionNumber % 100 >= 11 && positionNumber % 100 <= 13 ? suffixes[0] : suffixes[positionNumber % 10 <= 3 ? positionNumber % 10 : 0];
  return `${position}${suffix}`;
}

