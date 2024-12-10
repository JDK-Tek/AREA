/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** AreaDiscord1
*/

import { User } from "lucide-react";
import React, { useEffect, useState } from "react";
import HeaderBar from "../../components/Header/HeaderBar";
import Button from "../../components/Button";
import Notification from "../../components/Notification";
import fetchData from "../../utils/fetchData";

function DropdownBox({ options, selected, onSelect }) {
    const [isOpen, setIsOpen] = useState(false);

    const handleSelect = (option) => {
        onSelect(option);
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
    const [error, setError] = useState(null);
    const [success, setSuccess] = useState(null);
    const [time, setTime] = useState("");
    const [channel, setChannel] = useState("");
    const [message, setMessage] = useState("");
    const [unit, setUnit] = useState(options[0]);
    
    const services = {
        title: "Schedule sending of discord message",
        logo1: "/assets/services/discord.webp",
        logo2: "/assets/services/time.webp",
        users: 324434
    };

    const onSendArea = () => {
        if (!time || !channel || !message) {
            setError("Tous les champs sont obligatoires !");
            return;
        }

        const request = {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + sessionStorage.getItem("token")
            },
            body: JSON.stringify({
                action: {
                    service: "time",
                    name: "in",
                    spices: {
                        howmuch: time,
                        unit: unit
                    }
                },
                reaction: {
                    service: "discord",
                    name: "send",
                    spices: {
                        channel: channel,
                        message: message
                    }
                }
            })
        };

        console.log(request);

        fetchData("http://localhost:42000/api/area", request).then(({ success, data, error }) => {
            if (!success) {
                setError("Error while sending area: " + error);
            } else {
                setSuccess("Area has been sent successfully");
            }
        });
    };

    useEffect(() => {
        if (!sessionStorage.getItem("token")) {
            window.location.replace("/login");
        }
    }, []);

    return (
        <div>
            {error && <Notification msg={error} error={true} setError={setError}/>}
            {success && <Notification msg={success} error={false} setError={setSuccess}/>}
            <HeaderBar activeBackground={true}/>
            <AreaHeader services={services}/>

            <div className="flex justify-center items-center">
                <div>
                    <div className="flex items-center m-5">
                        <label className="text-[#7289da] mr-3 font-spartan font-bold text-lg">In</label>
                        <input
                            type="number"
                            id="time"
                            value={time}
                            onChange={(e) => setTime(e.target.value)}
                            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 mr-3"
                            placeholder="Time"
                        />
                        <DropdownBox options={options} selected={unit} onSelect={setUnit} />
                    </div>
                    <div className="flex items-center m-5">
                        <label className="text-[#7289da] mr-3 font-spartan font-bold text-lg">Channel</label>
                        <input
                            type="text"
                            id="channel"
                            value={channel}
                            onChange={(e) => setChannel(e.target.value)}
                            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 mr-3"
                            placeholder="Channel"
                        />
                    </div>
                    <div className="flex items-center m-5">
                        <label className="text-[#7289da] mr-3 font-spartan font-bold text-lg">Message</label>
                        <input
                            type="text"
                            id="message"
                            value={message}
                            onChange={(e) => setMessage(e.target.value)}
                            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 mr-3"
                            placeholder="Message"
                        />
                    </div>
                    <div className="flex items-center m-5">
                        <Button
                            text="Send"
                            styleClolor="bg-[#7289da] text-white"
                            onClick={onSendArea}
                        />
                    </div>
                </div>
            </div>
        </div>
    );
}
