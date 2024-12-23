/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** FindService
*/

import axios from "axios";
import { useEffect, useState } from "react";

import matchPattern from "../../utils/matchPattern";
import SearchInput from '../SearchInputBox'
import FeaturesKit from "./FeaturesKit";
import Notification from '../Notification'

export default function FindFeature({ dark }) {
    const [feature, setFeature] = useState({
        color: "",
        colorHover: "",
        feat: []
    });

    const [error, setError] = useState("");

    const [search, setSearch] = useState("");
    const [filteredFeature, setFilteredFeature] = useState({
        color: "",
        colorHover: "",
        feat: []
    });
    

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
        setFeature({
            color: "#05b348",
            colorHover: "#038a2b",
            feat: [
                { title: "When a new song is added to a playlist"},
                { title: "When a new song is played on a playlist"},
                { title: "When an followed artist releases a new album"},
                { title: "When an user likes one of your playlists"},
                { title: "When an user follows you"}
            ]
        })
    }, [setFeature, setError]);

    useEffect(() => {
        if (search === "") {
            setFilteredFeature(feature);
            return;
        }

        let fstmp = [];
        feature.feat.forEach((service) => {
            if (matchPattern(search, service.title)) {
                fstmp.push(service);
            }
        });
        const searchFeature = {
            color: feature.color,
            colorHover: feature.colorHover,
            feat: fstmp
        }
        setFilteredFeature(searchFeature);
    }, [search, setFilteredFeature, feature]);

    return (
        <div className="h-full flex flex-col justify-start">
            {error && <Notification error={true} msg={error} setError={setError}/>}

            <SearchInput
                placeholder={"Search for a service"}
                setText={setSearch}
    
                bgColor={mode.bgColor}
                txtColor={mode.txtColor}
                iconColor={mode.iconColor}
                borderColor={mode.borderColor}
            />
            
            <div className="mt-5 overflow-y-auto max-h-[calc(85vh-4rem-64px)]">
                <FeaturesKit features={filteredFeature} bgColor="bg-[#1d1d1d]"/>
            </div>
        </div>
    )
}
