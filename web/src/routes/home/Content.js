/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Content
*/

import AppletKit from "./../../components/Applet/AppletKit";
import ServiceKit from "./../../components/Service/ServiceKit";

import AppletData from "./../../data/AppletData";
import ServiceData from "./../../data/ServiceData";

export default function Content({ data }) {

    return (
        <div className="pb-14">
      
            <AppletKit
            title={"Get started with any Applet"}
            applets={AppletData}
            />
    
            <ServiceKit
            title={"or choose from 900+ services"}
            services={ServiceData}
            color={"text-chartpurple-200"}
            />
        
      </div>
    );
}