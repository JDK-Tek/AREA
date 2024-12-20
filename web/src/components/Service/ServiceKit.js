/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** ServiceKit
*/

import Service from "./Service"

export default function ServiceKit({ title, services, color = "text-black", gap, rounded, centered = true, bgColor = "" }) {

    return (
        <div className={`text-center p-5 ${bgColor}`}>
            <label className={`font-spartan ${color} text-[30px] font-bold`}>{title}</label>
            <div className={`mt-5 flex flex-wrap ${centered ? "justify-center" : ""} ${gap} p-5`}>
                {services.map((service, index) => (
                    <Service
                        key={index}
                        service={service}
                        rounded={rounded}
                    />
                ))}
            </div>
        </div>
    )
}
