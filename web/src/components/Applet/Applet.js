/*
** EPITECH PROJECT, 2024
** area
** File description:
** Applet
*/

import Icon from "./../../assets/icon.png"
import formatNumber from "./../../utils/FormatNumber";

export default function Applet({
        title,
        color,
        serviceAction,
        imageAction,
        serviceReaction,
        imageReaction,
        users,
        onClick
    }) {

    return (
        <div
            className={`cursor-pointer relative w-[300px] h-[325px] text-white rounded-3xl shadow-md p-4 transition-transform duration-200`}
            onClick={() => onClick()}
            onKeyDown={(event) => { if (event.key === "Enter") onClick() }}
            tabIndex={0}
            style={{
                backgroundColor: color
            }}
            onMouseEnter={(e) => { e.currentTarget.style.transform = "scale(1.05)"; }}
            onMouseLeave={(e) => { e.currentTarget.style.transform = "scale(1)"; }}
        >
            <div className="flex space-x-1">
                    <img
                        key={1}
                        src={`${process.env.REACT_APP_BACKEND_URL}${imageAction}`}
                        alt={serviceAction}
                        className="w-6 h-6"
                    />
                    <img
                        key={2}
                        src={`${process.env.REACT_APP_BACKEND_URL}${imageReaction}`}
                        alt={serviceReaction}
                        className="w-6 h-6"
                    />
            </div>

            <h1 className="font-spartan text-2xl font-bold mt-3 mb-3 text-left">{title}</h1>

            <div className="absolute bottom-6 w-full">
                <div className="flex items-center mt-2 space-x-2">
                    <img
                        src={`${process.env.REACT_APP_BACKEND_URL}${imageAction}`}
                        alt={serviceAction}
                        className="w-5 h-5"
                    />
                    <p className="text-xs font-bold">{serviceAction}</p>
                </div>

                {users && <div className="flex items-center mt-4">
                    <img
                        src={Icon}
                        alt="users"
                        className="w-5 h-5"
                    />
                    <p className="text-xs font-bold ml-2">{formatNumber(users)}</p>
                </div>}
            </div>
        </div>
    );
}
