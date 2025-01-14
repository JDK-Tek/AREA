/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** matchPattern
*/

export default function matchPattern(pattern, src) {
    if (typeof pattern !== 'string' || typeof src !== 'string') {
        console.warn('The pattern and the source must be strings');
    }
    const lsrc = src.toLowerCase();
    const lpattern = pattern.toLowerCase();

    return lsrc.includes(lpattern);
}
