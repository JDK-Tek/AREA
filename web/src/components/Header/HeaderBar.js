/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** HeaderBar
*/

import Button from "./../Button";

import Logo from './../../assets/fullLogo.png';
import listRoutes from "../../data/Routes";

export default function HeaderBar({ activeBackground = false }) {
    
    
    const isLogged = sessionStorage.getItem("token") === "" ? false : true;
    
    const dataRoutes = [
        {
            title: "Explore",
            link: listRoutes.explore
        },
        {
            title: "Stories",
            link: "/stories"
        },
        {
            title: (isLogged ? "My Applets" : "Login"),
            link: (isLogged ? listRoutes.myApplets : listRoutes.login)
        }
    ]
    
    return (
        <div className={`flex justify-between items-center p-3 relative z-10 ${activeBackground ?
            "bg-gradient-to-l from-chartpurple-200 via-chartpurple-300 to-chartgray-300" : ""
        }`}>
            <div>
                <img
                    className={`h-[50px] ${activeBackground ? "cursor-pointer" : ""}`}
                    src={Logo}
                    alt="logo"
                    onClick={() => (window.location.href = "/")}
                    onKeyDown={(event) => {
                        if (event.key === "Enter") {
                        window.location.href = "/";
                        }
                    }}
                    tabIndex={0}
                />

            </div>
            <div className="flex flex-wrap justify-end gap-7 items-center">
                {dataRoutes.map((route, index) => (
                    <label
                        key={index}
                        className="font-bold font-spartan text-white p-5 text-lg cursor-pointer hover:text-gray-200"
                        onClick={() => window.location.href = route.link}
                        onKeyDown={(event) => {if (event.key === "Enter") window.location.href = route.link}}
                        tabIndex={0}
                    >{route.title}</label>
                ))}
                {isLogged && (
                    <Button
                        text={"Create"}
                        redirect={false}
                        onClick={() => isLogged ? window.location.href = "/create" : "/login"}
                        styleClolor={"bg-white text-chartgray-300 hover:bg-gray-200"}
                    />
                )}
                <Button
                    text={isLogged ? "Logout" : "Get started"}
                    redirect={false}
                    onClick={() => {
                        if (isLogged) {
                            sessionStorage.removeItem("token");
                            window.location.href = "/";
                        } else {
                            window.location.href = "/register";
                        }
                    }}
                    styleClolor={"bg-white text-chartgray-300 hover:bg-gray-200"}
                />
            </div>
        </div>
    );
}
