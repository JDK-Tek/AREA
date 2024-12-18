/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** CreateArea
*/

import { useState } from "react";
import { Plus } from "lucide-react";

import SidePannel from "../../components/SidePannel"
import HeaderBar from "../../components/Header/HeaderBar"

import Button from "../../components/Button";

export default function App() {
    const [open, setOpen] = useState(false);

    return (
        <div className="relative">
            <HeaderBar activeBackground={true} />

            <div className="relative">
                <SidePannel title={"Chose an action"} setOpen={setOpen} open={open}/>

                <Button
                    text="Action"
                    styleClolor={`bg-chartpurple-200 text-white`}
                    onClick={() => setOpen(true)}
                    icon={<Plus />}
                />
            </div>
        </div>
    );
}
