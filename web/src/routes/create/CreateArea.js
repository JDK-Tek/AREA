/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** CreateArea
*/

import { useState } from "react";
import { Plus } from "lucide-react";
import { Trash2 } from "lucide-react";

import HeaderBar from "../../components/Header/HeaderBar"
import SidePannel from "../../components/SidePannel"
import Button from "../../components/Button";

function Triger({title, color, onClick}) {
    return (
        <div
            className="flex items-center justify-between w-full border-b-2 shadow-sm pl-4 pr-4 p-1"
            style={{backgroundColor: color}}
        >
            <div className="flex items-center">
                <label className="block text-2xl font-bold font-spartan text-white">{title}</label>
            </div>
            <Button
                styleClolor={`bg-chartpurple-200 hover:bg-chartpurple-100 text-white`}
                icon={<Trash2 />}
                onClick={onClick}
            />
        </div>
    )
}

export default function CreateArea() {
    const [open, setOpen] = useState(false);
    const [configAction, setConfigAction] = useState(true);
    const [area, setArea] = useState({
        name: "",
        actions: [],
        reactions: []
    });

    return (
        <div className="relative">
            <HeaderBar activeBackground={true} />

            <div className="relative">
                <SidePannel action={configAction} setOpen={setOpen} open={open} setArea={setArea}/>

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
                            {area.actions.map((action) =>
                                <Triger
                                    title={action.title}
                                    color={action.color}
                                    onClick={() => {
                                        setArea((prevArea) => ({
                                            ...prevArea,
                                            actions: prevArea.actions.filter((a) => a.title !== action.title),
                                        }));
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
                                    onClick={() => {
                                        setArea((prevArea) => ({
                                            ...prevArea,
                                            reactions: prevArea.reactions.filter((a) => a.title !== reaction.title),
                                        }));
                                    }}
                                />)
                            }
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
