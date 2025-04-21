export const formatPercentage = val =>
  isNaN(val) || val === null ? '0.00%' : `${Number(val).toFixed(2)}%`;
export const formatTime = ms =>
  ms === null || isNaN(ms) ? '-' : `${Number(ms).toFixed(2)}ms`;
