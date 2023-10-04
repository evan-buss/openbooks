/**
 * Applies a conditional attribute to an element. If the value is true, the attribute is applied.
 * @param attribute The attribute to apply. Will be prefixed with "data-".
 * @param value The value to check.
 */
export function conditionalAttribute(attribute: string, value: boolean) {
  return value ? { [`data-${attribute}`]: true } : {};
}
