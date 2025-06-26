declare module 'hijri-date' {
  export default class HijriDate extends Date {
    constructor(date?: Date | string | number);
    getFullYear(): number;
    getMonth(): number;
    getDate(): number;
  }
}

declare module 'hijri-date/lib/safe' {
  export default class HijriDate extends Date {
    constructor(date?: Date | string | number);
    getFullYear(): number;
    getMonth(): number;
    getDate(): number;
  }
}
