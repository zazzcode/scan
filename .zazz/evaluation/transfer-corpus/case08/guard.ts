export function accepts(value) {
  return typeof value === 'string' && /^[a-z0-9]+$/.test(value);
}
