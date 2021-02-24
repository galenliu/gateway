/**
 *
 * utils - Some utility export functions
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.*
 */
export declare function repeatChar(char: string, len: number): string;
export declare function padLeft(str: string, len: number): string;
export declare function padRight(str: string, len: number): string;
export declare function hexStr(num: number, len: number): string;
export declare function alignCenter(str: string, len: number): string;
/**
 * Prints formatted columns of lines. Lines is an array of lines.
 * Each line can be a single character, in which case a separator line
 * using that character is printed. Otherwise a line is expected to be an
 * array of fields.
 * The alignment argument is expected to be a string, one character per
 * column. '<' = left, '>' = right, '=' = center). If no alignment
 * character is found, then left alignment is assumed.
 */
export declare function printTable(alignment: string, lines: string[][]): void;
/**
 * Get the current time.
 *
 * @returns {String} The current time in the form YYYY-mm-ddTHH:MM:SS+00:00
 */
export declare function timestamp(): string;
//# sourceMappingURL=utils.d.ts.map