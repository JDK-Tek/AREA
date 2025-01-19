/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** FindService
*/

import { useEffect, useState } from "react";

import matchPattern from "../../utils/matchPattern";
import SearchInput from '../SearchInputBox'
import FeaturesKit from "./FeaturesKit";

export default function FindFeature({ dark, setFeature, service, action = null }) {
    const [search, setSearch] = useState("");
    const [filteredFeature, setFilteredFeature] = useState([]);

    const mode = dark ?
        {
            bgColor: "bg-chartgray-200",
            txtColor: "text-white placeholder-cahrtgray-100",
            iconColor: "text-cahrtgray-100",
            borderColor: "border-chartgray-100 focus:border-blue-500"
        } : 
        {
            bgColor: "bg-gray-50",
            txtColor: "text-gray-900 placeholder-gray-400",
            iconColor: "text-cahrtgray-100",
            borderColor: "border-gray-300 focus:border-blue-500"
        }

    useEffect(() => {

        let fstmp = [];
        if (!action || action === "action") {
            service.actions.forEach((a) => {
                if (search === "") {
                    fstmp.push(a);
                } else if (matchPattern(search, a.description)) {
                    fstmp.push(a);
                }
            });
        } 
        if (!action || action === "reaction") {
            service.reactions.forEach((r) => {
                if (search === "") {
                    fstmp.push(r);
                } else if (matchPattern(search, r.description)) {
                    fstmp.push(r);
                }
            });
        }

        setFilteredFeature(fstmp);
    }, [search, setFilteredFeature, service, action]);


    return (
        <div className="h-full flex flex-col justify-start">

            <SearchInput
                placeholder={"Search for a service"}
                setText={setSearch}
    
                bgColor={mode.bgColor}
                txtColor={mode.txtColor}
                iconColor={mode.iconColor}
                borderColor={mode.borderColor}
            />
            
            <div className="mt-5 overflow-y-auto max-h-[calc(85vh-4rem-64px)]">
                <FeaturesKit 
                    features={filteredFeature}
                    bgColor={"bg-[#1d1d1d]"}
                    setFeature={setFeature}
                    color={service.color}
                />
            </div>
        </div>
    )
}
