/*
** EPITECH PROJECT, 2025
** AREA
** File description:
** DatePicker
*/


export default function DatePicker({ value, setValue }) {

    return (
        <div className="relative">
            <input
                type="date"
                className="bg-gray-50 border border-gray-300 text-chartgray-200 text-sm rounded-lg
                    focus:outline-none focus:ring-2 focus:ring-chartpurple-100
                    block p-2.5"
                placeholder="Select date"
                value={value}
                onChange={(e) => setValue(e.target.value)}
            />
        </div>
    );
}

