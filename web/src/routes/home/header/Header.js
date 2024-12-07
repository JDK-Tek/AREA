/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Header
*/

import Title from "./Title"
import HeaderBar from "./HeaderBar"
import Button from "../../../components/Button"

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
        <div className="bg-gradient-to-l from-chartpurple-200 via-chartpurple-300 to-chartgray-200 mb-10">
            <HeaderBar routes={dataRoutes} />
            <div className="flex flex-col items-center justify-center pb-5">
                <Title />
                <Button
                    text={"Get started"}
                    redirect={true}
                    onClick={() => window.open("/register")}
                    styleClolor={"bg-white text-chartgray-200 hover:bg-gray-200 text-xl"}
                    />
            </div>
        </div>
    )
}
