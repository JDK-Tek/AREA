/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** SidePannel
*/

import { useEffect, useRef, useState } from "react";
import SearchInputBox from "./SearchInputBox";

export default function SidePannel({ title, open, setOpen }) {
    const panelRef = useRef(null);
    const [width, setWidth] = useState(400);
    const isResizing = useRef(false);

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

    useEffect(() => {
        const handleMouseMove = (event) => {
            if (isResizing.current) {
                const newWidth = Math.min(Math.max(400, event.clientX), window.innerWidth - 100);
                setWidth(newWidth);
            }
        };

        const handleMouseUp = () => {
            isResizing.current = false;
        };

        document.addEventListener("mousemove", handleMouseMove);
        document.addEventListener("mouseup", handleMouseUp);

        return () => {
            document.removeEventListener("mousemove", handleMouseMove);
            document.removeEventListener("mouseup", handleMouseUp);
        };
    }, []);

    return (
        <div
            ref={panelRef}
            style={{ width: `${width}px` }}
            className={`
                absolute left-0 h-[calc(100vh-4rem)] bg-chartgray-300 text-white z-10 shadow-lg
                transition-transform duration-300 ease-in-out ${open ? "translate-x-0" : "-translate-x-full"}`}
        >
            <div className="ml-5 mr-5 mt-3 flex flex-wrap justify-center items-center border-b-[1px] border-chartgray-200">
                <p className="p-3 font-spartan font-bold text-xl">{title}</p>
                <div
                    onMouseDown={() => (isResizing.current = true)}
                    style={{ cursor: "ew-resize" }}
                    className="absolute top-0 right-0 w-1 h-full bg-chartgray-200"
                />
            </div>
            <div className="p-5">
                <SearchInputBox placeholder={"Search a service ..."}/>
            </div>
        </div>
    );
}
