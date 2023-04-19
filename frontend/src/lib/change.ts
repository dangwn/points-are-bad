export const preventNegativeInputs = (e: React.KeyboardEvent<HTMLInputElement>): void => {
  if (e.key === '-' || e.key === '_') {
    e.preventDefault();
  };
};