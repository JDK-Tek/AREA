/*
** EPITECH PROJECT, 2025
** AREA
** File description:
** TextBox
*/

export default function TextBox({ value, setValue })
{
    return (
        <textarea
            id="message"
            rows="4"
            className="block p-2.5 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300
                    focus:outline-none focus:ring-2 focus:ring-chartpurple-100"
            placeholder={value}
            onChange={(e) => setValue(e.target.value)}
            value={value}
        />
    )
}
