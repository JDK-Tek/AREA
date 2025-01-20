/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Download
*/

import AndroidDownload from './../assets/android-download.png';
import AppleDownload from './../assets/apple-download.png';

export default function Download() {
    
    return (
        <div className="p-14">
            <img
                className="w-[200px] h-[60px] m-5 cursor-pointer"
                src={AndroidDownload}
                alt="Logo"
                onClick={() => window.open("https://play.google.com/store/apps?hl=fr")}
                onKeyDown={(event) => { if (event.key === "Enter") window.open("https://play.google.com/store/apps?hl=fr") }}
                tabIndex={0}
            />

            <img
                className="w-[200px] h-[60px] m-5 cursor-pointer"
                src={AppleDownload}
                alt="Logo"
                onClick={() => window.open("https://www.apple.com/fr/app-store/")}
                onKeyDown={(event) => { if (event.key === "Enter") window.open("https://www.apple.com/fr/app-store/") }}
                tabIndex={0}
            />
        </div>
    );
}