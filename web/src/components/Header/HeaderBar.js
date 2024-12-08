/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** HeaderBar
*/

import Button from "./../Button";

import Logo from './../../assets/fullLogo.png';

const dataRoutes = [
    {
        title: "Explore",
        link: "/about"
    },
    {
        title: "Stories",
        link: "/stories"
    },
    {
        title: "Login",
        link: "/login"
    }
]

export default function HeaderBar({ activeBackground = false }) {
    return (
        <div className={`flex justify-between items-center p-3 relative z-10 ${activeBackground ?
            "bg-gradient-to-l from-chartpurple-200 via-chartpurple-300 to-chartgray-300" : ""
        }`}>
            <div>
                <img
                    className={`h-[50px] ${activeBackground ? "cursor-pointer" : ""}`}
                    src={Logo}
                    alt="logo"
                    onClick={() => window.location.href = "/"}
                />
            </div>
            <div className="flex flex-wrap justify-end gap-7 items-center">
                {dataRoutes.map((route, index) => (
                    <label
                        key={index}
                        className="font-bold font-spartan text-white p-5 text-lg cursor-pointer hover:text-gray-200"
                        onClick={() => window.location.href = route.link}
                    >{route.title}</label>
                ))}
                <Button
                    text={"Get started"}
                    redirect={false}
                    onClick={() => window.location.href = "/register"}
                    styleClolor={"bg-white text-chartgray-300 hover:bg-gray-200"}
                />
            </div>
        </div>
    );
}
