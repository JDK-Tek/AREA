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
import { LRButton } from "./Button";

import Button from "./Button";
import { Undo2 } from 'lucide-react';

function handleOauth(service, setToken) {
    axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/oauth/${service}`)
        .then((response) => {
            const oauthWindow = window.open(response.data, "_blank");

            const handleMessage = (event) => {
                if (event.origin !== window.location.origin) {
                    return;
                }
                const code = event.data;
                if (code !== null) {
                    oauthWindow.close();
                    window.removeEventListener('message', handleMessage);
                    axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/oauth/${service}`, {
                        code: code
                    }, {
                        headers: {
                            "Content-Type": "application/json",
                            "Authorization": `Bearer ${sessionStorage.getItem("token")}`,
                        }
                    })
                    .then((response) => {
                        setToken(response.data.token);
                        window.location.reload();
                    })
                    .catch((error) => {
                        console.error('Error:', error);
                    });
                }
            };

            window.addEventListener('message', handleMessage);
        })
        .catch((error) => {
            console.error('Error:', error);
        });
}

export default function SidePannel({ action, open, setOpen, setArea, loggedServices, refresh, setToken }) {
    const panelRef = useRef(null);
    const [width, setWidth] = useState(540);
    const [service, setService] = useState(null);
    const [feature, setFeature] = useState(null);
    const isResizing = useRef(false);
    const [error, setError] = useState(false);
    const [errorMsg, setErrorMsg] = useState("");

    const [aboutjson, setAboutjson] = useState(null);
    const [content, setContent] = useState(null);

    useEffect(() => {
        const getAboutJson = async () => {
            axios.get(`${process.env.REACT_APP_BACKEND_URL}/about.json`, { headers: { "Content-Type": "application/json" } })
                .then((response) => {
                    setAboutjson(response.data);
                })
                .catch((error) => {
                    setError(true);
                    setErrorMsg("Error when trying to get about.json: " + error);
                });
        };
        getAboutJson();

    }, [aboutjson, setAboutjson, setError, setErrorMsg]);

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
    }, [setWidth, isResizing]);

    useEffect(() => {
        if (feature) {
            if (service.oauth && !loggedServices.some(s => s === service.name)) {
                 setContent(
                    <div className="flex flex-col items-center">
                        <LRButton
                            color={service.color}
                            img={`${process.env.REACT_APP_BACKEND_URL}${service.image}`}
                            text={"Connect to " + service.name}
                            handleClick={() => {
                                handleOauth(service.name, setToken);
                                refresh();
                            }}
                        />
                    </div>
                );
            } else {
                setContent(
                    <ServiceFeatureConfiguration
                        action={action}
                        feature={feature}
                        service={service}
                        setArea={setArea}
                        setError={setError}
                        setErrorMsg={setErrorMsg}
                        reset={() => {
                            setOpen(false);
                            setFeature(null);
                            setService(null);
                        }}
                    />
                );
            }
        } else if (service) {
            setContent(
                <FindFeature
                    dark={true}
                    setFeature={setFeature}
                    service={service}
                    action={action ? "action" : "reaction"}
                />
            );
        } else {
            setContent(
                <FindService
                    dark={true}
                    setService={setService}
                    setError={setError}
                    setErrorMsg={setErrorMsg}
                    aboutjson={aboutjson}
                    filtre={action ? "action" : "reaction"}
                />
            );
        }
    }, [
        feature,
        service,
        action,
        aboutjson,
        loggedServices,
        setArea,
        setOpen,
        refresh,
        setToken
    ]);

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
                {content}
            </div>
        </div>
    );
}
