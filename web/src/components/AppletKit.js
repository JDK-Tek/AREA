/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** AppletKit
*/
import Applet from "./Applet";

export default function AppletKit({ title, applets }) {
    return (
        <div className="text-center p-5">
            <label className="font-spartan text-[35px] font-bold">{title}</label>
            <div className="mt-5 flex flex-wrap justify-center gap-7 p-5">
                {applets.map((applet, index) => (
                    <Applet key={index} applet={applet} />
                ))}
            </div>
        </div>
    );
}
