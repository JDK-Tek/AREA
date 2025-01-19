/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Header
*/

import Title from "./Title"
import HeaderBar from "./HeaderBar"
import Button from "./../Button"

import BackgroundAsset from "./../../assets/effect.png"

export default function Header() {
    const isLogged = sessionStorage.getItem("token") === "" ? false : true;

    return (
        <div className="relative bg-gradient-to-l from-chartpurple-200 via-chartpurple-300 to-chartgray-300 mb-10 overflow-hidden">
            <div className="absolute inset-0 z-0">
                <img
                    src={BackgroundAsset}
                    alt="Background pattern"
                    className="absolute top-[-300px] left-[-200px] w-[700px] rotate-[165deg]"
                />
                <img
                    src={BackgroundAsset}
                    alt="Background wave"
                    className="absolute top-[-400px] right-[-400px] w-[650px] rotate-[165deg]"
                />
            </div>
    
            <HeaderBar />

            <div className="flex flex-col items-center justify-center pb-5 relative z-10">
                <Title />
                {!isLogged && <Button
                    text={"Get started"}
                    redirect={true}
                    onClick={() => window.open("/register")}
                    styleClolor={"bg-white text-chartgray-300 hover:bg-gray-200 text-xl"}
                />}
            </div>
        </div>
    );
}
