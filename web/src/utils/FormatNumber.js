/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** FormatNumber
*/

const formatting = [
    {
        value: 1000000000,
        extension: 'b'
    },
    {
        value: 1000000,
        extension: 'm'
    },
    {
        value: 1000,
        extension: 'k'
    }
];

export default function formatNumber(num){
    for (const format of formatting) {
        if (num >= format.value) {
            return (num / format.value).toFixed(1).replace('.0', '') + format.extension;
        }
    }
    return num.toString();
}
