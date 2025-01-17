/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** SidePannel
*/

import axios from "axios";
import { useEffect, useRef, useState } from "react";

import FindService from "./Service/FindService";
import FindFeature from "./Feature/FindFeature";
import ServiceFeatureConfiguration from "./ServiceFeatureConfiguration";
import Notification from "./Notification";

import Button from "./Button";
import { Undo2 } from 'lucide-react';

export default function SidePannel({ action, open, setOpen, setArea }) {
    const panelRef = useRef(null);
    const [width, setWidth] = useState(540);
    const [service, setService] = useState(null);
    const [feature, setFeature] = useState(null);
    const isResizing = useRef(false);
    const [error, setError] = useState(false);
    const [errorMsg, setErrorMsg] = useState("");

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
                const newWidth = Math.min(Math.max(330, event.clientX), window.innerWidth / 2 - 100);
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
                absolute left-0 h-[calc(100.4vh-6rem)] bg-chartgray-300 text-white z-10 shadow-lg
                transition-transform duration-300 ease-in-out ${open ? "translate-x-0" : "-translate-x-full"}`}
        >
            {error && <Notification error={true} msg={errorMsg} setError={setError}/>}
            <div className="ml-5 mr-5 mt-3 flex flex-wrap justify-center items-center border-b-[1px] border-chartgray-200">
                {service &&
                    <div className="absolute left-0">
                        <Button
                            onClick={() => {
                                if (feature) setFeature(null)
                                else setService(null)
                            }}
                            icon={<Undo2 />}
                            styleClolor={"bg-chartpurple-200 text-white hover:bg-chartpurple-100 text-2xl"}
                        />
                    </div>
                }
                <p className="p-3 font-spartan font-bold text-2xl">{action ? "Configurate an action" : "Configurate a reaction"}</p>
                <div
                    onMouseDown={() => (isResizing.current = true)}
                    style={{ cursor: "ew-resize" }}
                    className="absolute top-0 right-0 w-1 h-full bg-chartgray-200"
                />
            </div>
            <div className="p-5" style={{ height: 'calc(95vh - 4rem - 64px)' }}>
                {
                feature ?
                    <ServiceFeatureConfiguration
                        action={action}
                        color={service.color}
                        setError={setError}
                        setErrorMsg={setErrorMsg}
                        setArea={setArea}
                        reset={() => {
                            setOpen(false);
                            setFeature(null);
                            setService(null);
                        }}
                    />
                    :
                service ?
                    <FindFeature
                        dark={true}
                        setFeature={setFeature}
                        color={service.color}
                    />
                    :
                    <FindService
                        dark={true}
                        setService={setService}
                        setError={setError}
                        setErrorMsg={setErrorMsg}
                    />
                }
            </div>
        </div>
    );
}
