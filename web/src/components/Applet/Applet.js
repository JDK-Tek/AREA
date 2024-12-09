/*
** EPITECH PROJECT, 2024
** area
** File description:
** Applet
*/

import Icon from "./../../assets/icon.png"
import formatNumber from "./../../utils/FormatNumber";

export default function Applet({ applet }) {

    return (
        <div
            className={`cursor-pointer relative w-[300px] h-[325px] text-white rounded-3xl shadow-md p-4`}
            onClick={() => window.location.href = applet.link}
            style={{
                backgroundColor: applet.service.color.normal
            }}
            onMouseEnter={(e) => e.currentTarget.style.backgroundColor = applet.service.color.hover}
            onMouseLeave={(e) => e.currentTarget.style.backgroundColor = applet.service.color.normal}
        >
            <div className="flex space-x-1">
                    <img
                        key={1}
                        src={applet.service.logo}
                        alt={applet.service.name}
                        className="w-6 h-6"
                    />
                    <img
                        key={2}
                        src={applet.service.logopartner}
                        alt={applet.service.name}
                        className="w-6 h-6"
                    />
            </div>

            <h1 className="font-spartan text-2xl font-bold mt-3 mb-3 text-left">{applet.name}</h1>

            <div className="absolute bottom-6 w-full">
                <div className="flex items-center mt-2 space-x-2">
                    <img
                        src={applet.service.logo}
                        alt={applet.service.name}
                        className="w-5 h-5"
                    />
                    <p className="text-xs font-bold">{applet.service.name}</p>
                </div>

                <div className="flex items-center mt-4">
                    <img
                        src={Icon}
                        alt="users"
                        className="w-5 h-5"
                    />
                    <p className="text-xs font-bold ml-2">{formatNumber(applet.users)}</p>
                </div>
            </div>
        </div>
    );
}
