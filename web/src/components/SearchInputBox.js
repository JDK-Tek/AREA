/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** SearchInputBox
*/

export default function SearchInputBox({
        placeholder,
        setText,
        bgColor,
        iconColor,
        txtColor,
        borderColor
    }) {


    return (
        <div className="w-full pl-3 pr-3">
            <label
                className="mb-2 text-sm font-medium sr-only text-white"
            >
                Search
            </label>
            <div className="relative">
                <div className="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
                    <svg 
                        className={`w-4 h-4 ${iconColor}`}
                        aria-hidden="true"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 20 20"
                    >
                        <path
                            stroke="currentColor"
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth="2"
                            d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"
                        />
                    </svg>
                </div>
                <input
                    id="default-search"
                    className={`placeholder:select-none block w-full p-4 ps-10 text-sm ${txtColor} border ${borderColor} rounded-lg ${bgColor}`}
                    placeholder={placeholder}
                    onChange={(e) => setText(e.target.value)}
                />
            </div>  
        </div>
    )
}