"use strict";
/**
 *
 * utils - Some utility export functions
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.*
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.timestamp = exports.printTable = exports.alignCenter = exports.hexStr = exports.padRight = exports.padLeft = exports.repeatChar = void 0;
function repeatChar(char, len) {
    if (len <= 0) {
        return '';
    }
    return new Array(len + 1).join(char);
}
exports.repeatChar = repeatChar;
function padLeft(str, len) {
    return (repeatChar(' ', len) + str).slice(-len);
}
exports.padLeft = padLeft;
function padRight(str, len) {
    return (str + repeatChar(' ', len)).slice(0, len);
}
exports.padRight = padRight;
function hexStr(num, len) {
    return (repeatChar('0', len) + num.toString(16)).slice(-len);
}
exports.hexStr = hexStr;
function alignCenter(str, len) {
    if (str.length >= len) {
        return str.slice(0, len);
    }
    const leftSpace = (len - str.length) / 2;
    return padRight(padLeft(str, str.length + leftSpace), len);
}
exports.alignCenter = alignCenter;
/**
 * Prints formatted columns of lines. Lines is an array of lines.
 * Each line can be a single character, in which case a separator line
 * using that character is printed. Otherwise a line is expected to be an
 * array of fields.
 * The alignment argument is expected to be a string, one character per
 * column. '<' = left, '>' = right, '=' = center). If no alignment
 * character is found, then left alignment is assumed.
 */
function printTable(alignment, lines) {
    const width = [];
    let colWidth;
    let idx;
    // Take a pass through the data and figure out the width for each column.
    for (const line of lines) {
        if (typeof line !== 'string') {
            for (idx in line) {
                colWidth = line[idx].length;
                if (typeof width[idx] === 'undefined' || colWidth > width[idx]) {
                    width[idx] = colWidth;
                }
            }
        }
    }
    // Now that we know how wide each column is, go and print them.
    for (const line of lines) {
        let lineStr = '';
        for (const idx in width) {
            if (parseInt(idx) > 0) {
                lineStr += ' ';
            }
            colWidth = width[idx];
            if (typeof line === 'string') {
                lineStr += repeatChar(line[0], colWidth);
            }
            else {
                const align = alignment[idx];
                let field = line[idx];
                if (typeof field === 'undefined') {
                    field = '';
                }
                if (align === '>') {
                    lineStr += padLeft(field, colWidth);
                }
                else if (align === '=') {
                    lineStr += alignCenter(field, colWidth);
                }
                else {
                    lineStr += padRight(field, colWidth);
                }
            }
        }
        console.log(lineStr);
    }
}
exports.printTable = printTable;
/**
 * Get the current time.
 *
 * @returns {String} The current time in the form YYYY-mm-ddTHH:MM:SS+00:00
 */
function timestamp() {
    const date = new Date().toISOString();
    return date.replace(/\.\d{3}Z/, '+00:00');
}
exports.timestamp = timestamp;
//# sourceMappingURL=utils.js.map