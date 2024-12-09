/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Notification
*/

import { TriangleAlert, Check, X } from "lucide-react";

export default function Notification({ error, msg, setError }) {
    return (
        <div className="fixed inset-0 flex items-start justify-center z-[100] pointer-events-none">
            <div
                role="alert"
                className={`relative mt-10 max-w-[500px] w-full flex flex-col p-3 text-sm text-white ${
                    error ? "bg-red-500" : "bg-green-600"
                } rounded-md shadow-lg pointer-events-auto`}
            >
                <div className="flex items-start">
                    {error ? ( <TriangleAlert /> ) : ( <Check /> )}
                    <p className="ml-3 text-base font-bold">{error ? "Error" : "Success"}</p>
                </div>
                <p className="ml-7 mt-2">{msg}</p>
                <button
                    onClick={() => setError(false)}
                    className="absolute top-2 right-2 flex items-center justify-center transition-all w-8 h-8 rounded-md text-white hover:bg-white/10 active:bg-white/20"
                    type="button"
                >
                    <X />
                </button>
            </div>
        </div>
    );
}

