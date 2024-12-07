/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** NotFound
*/

import HeaderBar from "./../../components/Header/HeaderBar";
import Logo from "./../../assets/logo.png";

export default function NotFound() {
    return (
        <div>
            <HeaderBar activeBackground={true}/>
            <div className="min-h-screen flex flex-col items-center justify-center">
                <img src={Logo} alt="logo" className="mb-6" />
                <h1 className="text-5xl font-bold font-spartan text-chartpurple-200">404 - Not Found</h1>
            </div>
        </div>
    );
}
