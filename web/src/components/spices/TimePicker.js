/*
** EPITECH PROJECT, 2025
** AREA
** File description:
** TimePicker
*/

export default function TimePicker({ value, setValue }) {
    return (
        <div className="relative">
            <input
                type="time"
                id="time"
                className="bg-gray-50 border leading-none border-gray-300 text-chartgray-200 text-sm rounded-lg
                    focus:ring-blue-500 focus:border-blue-500
                    block p-2.5"
                min="09:00"
                max="18:00"
                value={value}
                onChange={(e) => setValue(e.target.value)}
                required
            />
        </div>
    );
}

