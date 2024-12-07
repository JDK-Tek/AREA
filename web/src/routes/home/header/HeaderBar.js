/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** HeaderBar
*/

import Button from "../../../components/Button";

import Logo from './../../../assets/fullLogo.png';

export default function HeaderBar({ routes }) {
    return (
        <div className="flex justify-between items-center p-3 relative z-10">
            <div>
                <img src={Logo} alt="logo" className="h-[50px]"/>
            </div>
            <div className="flex flex-wrap justify-end gap-7 items-center">
                {routes.map((route, index) => (
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
