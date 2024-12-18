/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** SidePannel
*/

import { useEffect, useRef } from "react";

export default function SidePannel({ title, open, setOpen }) {
    const panelRef = useRef(null);

    
    useEffect(() => {
    
        const handleClickOutside = (event) => {
            if (panelRef.current && !panelRef.current.contains(event.target)) {
                setOpen(false);
            }
        };

        if (open) {
            document.addEventListener("mousedown", handleClickOutside);
        }

        return () => {
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, [open, setOpen]);

    return (
        <div
            ref={panelRef}
            className={`
                absolute left-0 w-[400px] h-[calc(100vh-4rem)] bg-chartgray-300 text-white z-10 shadow-lg
                transition-transform duration-300 ease-in-out ${open ? "translate-x-0" : "-translate-x-full"}`}
        >
            <p className="p-4">{title}</p>
        </div>
    );
}
