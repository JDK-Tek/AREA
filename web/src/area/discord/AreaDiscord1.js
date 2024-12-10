/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** AreaDiscord1
*/

import { User } from "lucide-react"
import React, { useState } from "react";
import HeaderBar from "../../components/Header/HeaderBar";
import Button from "../../components/Button";

function DropdownBox({options}) {
    const [selected, setSelected] = useState(options[0]);
    const [isOpen, setIsOpen] = useState(false);

    const handleSelect = (option) => {
        setSelected(option);
        setIsOpen(false);
    };

    return (
        <div className="relative w-64">
            <button
                className="w-full bg-gray-200 text-gray-700 px-4 py-2 rounded-md shadow-md focus:outline-none focus:ring-2 focus:ring-blue-500 flex justify-between items-center"
                onClick={() => setIsOpen(!isOpen)}
            >
                <span>{selected}</span>
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



function AreaHeader({ services }) {
    return (
        <div className="flex justify-center items-center">
            <div className="w-full h-[500px] bg-[#7289da] m-10 rounded-xl p-10 shadow-md flex flex-col items-center justify-center">
                <div className="max-w-[350px] flex flex-col">
                    <div className="flex space-x-1 mb-3">
                        <img 
                            alt="first service"
                            src={services.logo1}
                            className="w-[65px]"
                        />
                        <img 
                            alt="second service"
                            src={services.logo2}
                            className="w-[65px]"
                        />
                    </div>
                    <label className="text-white font-spartan font-bold text-[40px] leading-[45px] mb-20">{services.title}</label>

                    <div className="flex">
                        <User color="white"/>
                        <label className="text-white font-spartan font-bold text-[20px] ml-2">{services.users}</label>
                    </div>
                </div>
            </div>
        </div>
    );
}



export default function AreaDiscord1() {
    const options = ["seconds", "minutes", "hours", "days", "weeks", "months", "years"];
    const services = {
        title: "Schedule sending of discord message",
        logo1: "/assets/services/discord.webp",
        logo2: "/assets/services/time.webp",
        users: 324434
    };

    return (
        <div>
            <HeaderBar activeBackground={true}/>
            <AreaHeader services={services}/>

            <div className="flex justify-center items-center">
                <div>
                    <div className="flex items-center m-5">
                        <label className="text-[#7289da] mr-3 font-spartan font-bold text-lg">In</label>
                        <input
                            type="number"
                            id="time"
                            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 mr-3 " placeholder="Time"
                            />
                        <DropdownBox options={options}/>
                    </div>
                    <div className="flex items-center m-5">
                        <label className="text-[#7289da] mr-3 font-spartan font-bold text-lg">Send</label>
                        <input
                            type="text"
                            id="time"
                            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 mr-3 " placeholder="Message"
                            />
                        <Button text="Send" styleClolor="bg-[#7289da] text-white"/>
                    </div>
                </div>
            </div>

        </div>
    );
}
