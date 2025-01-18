/*
** EPITECH PROJECT, 2025
** AREA
** File description:
** MobileClientDownload
*/

import React from 'react';
import Apk from '../../apk/area_jepgo.apk';


export default function MobileClientDownload() {
    return (
        <div>
            <div className="flex justify-center items-center mt-8">
                <a
                    href={Apk}
                    download="area_jepgo"
                    target="_blank"
                    rel="noreferrer"
                >
                    <button>Download APK file</button>
                </a>             
            </div>
        </div>
    )
}