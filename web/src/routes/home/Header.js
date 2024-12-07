/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Header
*/

import HeaderBar from "./HeaderBar"

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
        <div className="bg-gradient-to-l from-chartpurple-200 via-gray-800 to-gray-900">
            <HeaderBar routes={dataRoutes} />
        </div>
    )
}
