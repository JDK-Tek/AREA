/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Home
*/

import { useState } from "react";

import HeaderBar from "./../../components/Header/HeaderBar";
import Footer from "../home/Footer";
import Content from './Content';
import Notification from "../../components/Notification";

import BackgroundAsset from "./../../assets/effect.png"


export default function Login( {setToken} ) {
    const [error, setError] = useState(null);

    return (
        <div className="bg-gradient-to-tl from-chartpurple-200 via-chartpurple-300 to-chartgray-300 overflow-x-hidden relative h-screen w-full">
            <div className="absolute inset-0 z-0 overflow-hidden">
                <img
                    src={BackgroundAsset}
                    alt="Background pattern"
                    className="absolute top-[-300px] left-[-200px] w-[700px] rotate-[165deg]"
                />
                <img
                    src={BackgroundAsset}
                    alt="Background wave"
                    className="absolute top-[300px] right-[-600px] w-[1000px] rotate-[250deg]"
                />
            </div>
    
            <div className="relative z-10">
                <HeaderBar />
                {error && <Notification error={true} msg={error} setError={setError} />}
                <Content setToken={setToken} setError={setError} />
                <Footer />
            </div> 
        </div>
    );
}
