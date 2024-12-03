/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** ServiceKit
*/

import Service from "./Service"

export default function ServiceKit({ title, services, color = "text-black" }) {

    return (
        <div className="text-center p-5">
            <label className={`font-spartan ${color} text-[30px] font-bold`}>{title}</label>
            <div className="mt-5 flex flex-wrap justify-center gap-7 p-5">
                {services.map((service, index) => (
                    <Service
                        key={index}
                        service={service}
                    />
                ))}
            </div>
        </div>
    )
}
