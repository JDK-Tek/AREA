/*
** EPITECH PROJECT, 2025
** AREA [WSL : Ubuntu]
** File description:
** connected
*/

import React from "react";
import { useEffect } from "react";
import axios from "axios";

export default function Connected({ setToken }) {
    useEffect(() => {
        const urlParams = new URLSearchParams(window.location.search);
        const code = urlParams.get('code');
        
        const sendCode = async () => {
            try {
                window.opener.postMessage(code, window.opener.location.origin);
            } catch (error) {
                console.error('Error:', error);
            }
        };
        
        sendCode();
    }, []);

    return (
        <div className="bg-gradient-to-br from-zinc-900 via-indigo-900 to-violet-900 animate-gradient bg-300% h-screen justify-center items-center flex">
            <p className="text-white text-4xl sm:text-5xl md:text-6xl lg:text-7xl font-spartan font-bold">
                Connection in progress..
            </p>
        </div>
    );
}