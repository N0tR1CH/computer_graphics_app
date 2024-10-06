export type Shape = {
  name: string;
  x: number;
  y: number;
  height: number;
  width: number;
  /* Triangle needs base */
  base: number;
  /* Ellipse needs:
   *   - radius1
   *   - radius2
   *   - base
   *   - height
   */
  radius1: number;
  radius2: number;
  rotation: number;
  /* Straight line
  *   - x1
  *   - y1
  */
  x1: number;
  y1: number;
  /* Text
  *   - text
  */
  text: string;
};
