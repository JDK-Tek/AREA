/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** SearchInputBox
*/

export default function SearchInputBox({ placeholder, value, onChange }) {
    return (
        <div class="w-full pl-3 pr-3">
            <label for="default-search" class="mb-2 text-sm font-medium sr-only text-white">Search</label>
            <div class="relative">
                <div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
                    <svg class="w-4 h-4 text-chartgray-100" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20">
                        <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"/>
                    </svg>
                </div>
                <input type="search" id="default-search" class="placeholder:select-none block w-full p-4 ps-10 text-sm text-white border border-chartgray-100 rounded-lg bg-chartgray-200" placeholder={placeholder} />
                <button type="submit" class="select-none text-white absolute end-2.5 bottom-2.5 bg-chartpurple-100 hover:bg-chartpurple-200 active:ring-2 active:outline-none active:border-chartpurple-200 font-medium rounded-lg text-sm px-4 py-2">Search</button>
            </div>  
        </div>
    )
}