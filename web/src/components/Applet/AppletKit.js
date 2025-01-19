/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** AppletKit
*/
import Applet from "./Applet";

export default function AppletKit({
        title,
        applets,
        color = "text-black",
        onClick
    }) {

    return (
        <div className="text-center p-5">
            <label className={`font-spartan ${color} text-[35px] font-bold`}>{title}</label>
            <div className="mt-5 flex flex-wrap justify-center gap-7 p-5">
                {applets.map((applet, index) => (
                    <Applet
                        key={index}
                        title={applet.name}
                        color={applet.action.color}
                        serviceAction={applet.action.service}
                        imageAction={applet.action.image}
                        serviceReaction={applet.reaction.service}
                        imageReaction={applet.reaction.image}
                        users={applet.users}
                        onClick={onClick}
                    />
                ))}
            </div>
        </div>
    );
}
