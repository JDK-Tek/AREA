/*
** EPITECH PROJECT, 2025
** AREA
** File description:
** NumberInputBox
*/

export default function NumberInputBox({ value, setValue }) {
    return (
        <input
            type="number"
            id="number-input"
            aria-describedby="helper-text-explanation"
            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg
                    focus:outline-none focus:ring-2 focus:ring-chartpurple-100
                    block w-full p-2.5"
            placeholder={value}
            onChange={(e) => setValue(e.target.value)}
        />
    )
}