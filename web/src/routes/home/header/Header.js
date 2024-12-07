/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Header
*/

import Title from "./Title"
import HeaderBar from "./HeaderBar"
import Button from "../../../components/Button"

import BeckgroundAssets from "../../../assets/effect.png"

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


export default function Header() {

    return (
        <div className="relative bg-gradient-to-l from-chartpurple-200 via-chartpurple-300 to-chartgray-300 mb-10 overflow-hidden">
            {/* Background Images */}
            <div className="absolute inset-0 z-0">
                <img
                    src={BeckgroundAssets}
                    alt="Background pattern"
                    className="absolute top-[-300px] left-[-200px] w-[700px] rotate-[165deg]"
                />
                <img
                    src={BeckgroundAssets}
                    alt="Background wave"
                    className="absolute top-[-400px] right-[-400px] w-[650px] rotate-[165deg]"
                />
            </div>
    
            {/* HeaderBar */}
            <HeaderBar routes={dataRoutes} />

            {/* Main Content */}
            <div className="flex flex-col items-center justify-center pb-5 relative z-10">
                <Title />
                <Button
                    text={"Get started"}
                    redirect={true}
                    onClick={() => window.open("/register")}
                    styleClolor={"bg-white text-chartgray-300 hover:bg-gray-200 text-xl"}
                />
            </div>
        </div>
    );
}
