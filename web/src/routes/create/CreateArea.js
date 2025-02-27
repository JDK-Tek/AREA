/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** CreateArea
*/

import axios from "axios";
import { useEffect, useState } from "react";
import { Plus, PencilLine, Trash2 } from "lucide-react";

import HeaderBar from "../../components/Header/HeaderBar"
import SidePannel from "../../components/SidePannel"
import Button from "../../components/Button";
import InputBox from "../../components/spices/InputBox";
import Notification from "../../components/Notification";

function Triger({title, color, spices, onClickTrash, onClickPencil}) {
    return (
        <div
            className="flex items-center justify-between w-full border-b-2 shadow-sm pl-4 pr-4 p-1"
            style={{backgroundColor: color}}
        >
            <div className="flex items-center">
                <label className="block text-2xl font-bold font-spartan text-white">{title}</label>
            </div>
            <div>
                <Button
                    styleClolor={`bg-chartpurple-200 hover:bg-chartpurple-100 text-white`}
                    icon={<PencilLine />}
                    onClick={onClickPencil}
                />
                <Button
                    styleClolor={`bg-chartpurple-200 hover:bg-chartpurple-100 text-white`}
                    icon={<Trash2 />}
                    onClick={onClickTrash}
                />
            </div>
        </div>
    )
}

export default function CreateArea({setToken}) {
    const [open, setOpen] = useState(false);
    const [configAction, setConfigAction] = useState(true);
    const [name, setName] = useState("");
    const [error, setError] = useState(false);
    const [success, setSuccess] = useState(false);

    const [loggedServices, setLoggedServices] = useState([]);

    const defaultArea = {
        name: "",
        actions: [],
        reactions: []
    };
    
    const [area, setArea] = useState(sessionStorage.getItem("area") === null ? defaultArea : JSON.parse(sessionStorage.getItem("area")));
    sessionStorage.setItem("area", JSON.stringify(area));
    
    console.log(sessionStorage.getItem("token"));
    const checkConnection = () => {
        axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/doctor`, {
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${sessionStorage.getItem("token")}`,
            }
        })
        .then((res) => {
            if (res.data.authentificated && res.data.oauths) {
                setLoggedServices(res.data.oauths);
            } else {
                console.log("Not authentificated");
                console.log(res.data);
                // window.location.href = "/login";
                // sessionStorage.removeItem("token");
            }
        })
        .catch((err) => {
            setError("Impossible to check the authentification: " + err.data);
        });
    }

    useEffect(() => {
        checkConnection();
    }, [setLoggedServices]);


    return (
        <div className="relative">
            <HeaderBar activeBackground={true} />
            {error && <Notification error={error} setError={setError} msg={error} />}
            {success && <Notification success={success} setError={setSuccess} msg={success} />}

            <SidePannel 
                action={configAction} 
                setOpen={setOpen}
                open={open}
                setArea={setArea}
                loggedServices={loggedServices}
                refresh={checkConnection}
                setToken={setToken}
            />
            <div className="relative">
                <label 
                    className="block text-4xl font-bold font-spartan text-chartgray-300 text-center p-5 mt-10"
                >Create a new AREA</label>


                <div className="p-10 overflow-y-auto max-h-[calc(100vh-4rem-64px)]">
                    <div className="mb-10">
                        <div className="flex items-center w-full border-b-2 shadow-sm">
                            <Button
                                styleClolor={`bg-chartpurple-200 hover:bg-chartpurple-100 text-white`}
                                onClick={() => {
                                    setOpen(true)
                                    setConfigAction(true)
                                }}
                                icon={<Plus />}
                                />
                            <label className="block text-2xl font-bold font-spartan text-chartpurple-200">Configurate an action</label>
                        </div>
                        <div className="p-5">
                            {area.actions.map((action, index) =>
                                <Triger
                                    key={index}
                                    title={action.title}
                                    color={action.color}
                                    spices={action.spices}
                                    onClickTrash={() => {
                                        setArea((prevArea) => ({
                                            ...prevArea,
                                            actions: []
                                        }));
                                    }}
                                    onClickPencil={() => {
                                        setOpen(true);
                                        setConfigAction(true);
                                    }}
                                />)
                            }
                        </div>
                    </div>
                    <div className="">
                        <div className="flex items-center w-full border-b-2 shadow-sm">
                            <Button
                                styleClolor={`bg-chartgray-200 hover:bg-chartgray-100 text-white`}
                                onClick={() => {
                                    setOpen(true)
                                    setConfigAction(false)
                                }}
                                icon={<Plus />}
                                />
                            <label className="block text-2xl font-bold font-spartan text-chartgray-200">Configurate a reaction</label>
                        </div>
                        <div className="p-5">
                            {area.reactions.map((reaction) =>
                                <Triger
                                    title={reaction.title}
                                    color={reaction.color}
                                    spices={reaction.spices}
                                    onClickTrash={() => {
                                        setArea(defaultArea);
                                    }}
                                    onClickPencil={() => {
                                        setOpen(true);
                                        setConfigAction(false);
                                    }}
                                />)
                            }
                        </div>
                    </div>
                <div className="flex justify-center mt-10">
                    <Button
                        text="Create the new AREA"
                        styleClolor="bg-chartpurple-200 hover:bg-chartpurple-100 text-white"
                        onClick={() => {

                            if (area.actions.length === 0) {
                                setError("Missing at least one action");
                                return;
                            }

                            if (area.reactions.length === 0) {
                                setError("Missing at least one reaction");
                                return;
                            }

                            const token = sessionStorage.getItem("token");
                            const body = {
                                action: {
                                    service: area.actions[0].service,
                                    name: area.actions[0].name,
                                    spices: area.actions[0].spices,
                                },
                                reaction: {
                                    service: area.reactions[0].service,
                                    name: area.reactions[0].name,
                                    spices: area.reactions[0].spices,
                                }
                            };
                            const header = {
                                    "Content-Type": "application/json",
                                    "Authorization": `Bearer ${token}`,
                            };

                            axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/area`, body, {
                                headers: header
                            })
                            .then((res) => {
                                setSuccess("AREA created with success");
                                setArea(defaultArea);
                                sessionStorage.setItem("area", JSON.stringify(defaultArea));
                            })
                            .catch((err) => {
                                setError("An error occured while creating the AREA: " + err.data);
                            });
                        }}
                    />
                </div>
                </div>
            </div>
        </div>
    );
}
