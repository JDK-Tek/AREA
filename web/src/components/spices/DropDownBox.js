/*
** EPITECH PROJECT, 2025
** AREA
** File description:
** DropDownBox
*/

import { useState } from "react";

export default function DropdownBox({ options, selected, onSelect }) {
    const [isOpen, setIsOpen] = useState(false);

    const handleSelect = (option) => {
        onSelect(option);
        setIsOpen(false);
    };

    return (
        <div className="relative w-64">
            <button
                className="bg-gray-200 text-chartgray-200 px-4 py-2 rounded-md shadow-md
                        focus:outline-none focus:ring-2 focus:ring-chartpurple-100
                        flex justify-between items-center"
                onClick={() => setIsOpen(!isOpen)}
            >
                <span className="mr-5">{selected || "Select an option"}</span>
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-5 w-5 text-gray-500"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M19 9l-7 7-7-7"
                    />
                </svg>
            </button>

            {isOpen && (
                <div className="absolute mt-2 w-full bg-white border border-gray-300 rounded-md shadow-lg z-10 max-h-20 overflow-y-auto">
                    {options.map((option, index) => (
                        <div
                            key={index}
                            onClick={() => handleSelect(option)}
                            onKeyDown={(event) => {
                                if (event.key === "Enter") handleSelect(option);
                            }}
                            tabIndex={0}
                            className="px-4 py-2 text-gray-700 hover:bg-blue-500 hover:text-white cursor-pointer"
                        >
                            {option}
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
