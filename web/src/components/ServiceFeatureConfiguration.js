/*
** EPITECH PROJECT, 2025
** AREA
** File description:
** ServiceFeatureConfiguration
*/

import { useState } from "react";

import Button from "./Button";

import TextBox from "./spices/TextBox";
import DropdownBox from "./spices/DropDownBox";
import NumberInputBox from "./spices/NumberInputBox";
import DatePicker from "./spices/DatePicker";
import TimePicker from "./spices/TimePicker";
import InputBox from "./spices/InputBox";

function checkRequest(config, request) {
    const spices = config.spices;
    const missingParameters = [];

    spices.forEach((spice) => {
        if (!request[spice.name]) {
            if (missingParameters.length === 0) {
                missingParameters.push("Missing parameters: ");
            }
            missingParameters.push("'" + spice.title + "', ");
        }
    });
    return missingParameters;
}

export default function ServiceFeatureConfiguration({ action, feature, service, setArea, setError, setErrorMsg, reset }) {
    const [request, setRequest] = useState({});

    const handleValueChange = (name, value) => {
        setRequest((prevRequest) => ({
            ...prevRequest,
            [name]: value,
        }));
    };

    const renderInput = (spice) => {
        const value = request[spice.name] || (spice.type === "dropdown" ? spice.extra[0] : "");
    
        switch (spice.type) {
            case "text":
                return (
                    <TextBox
                        key={spice.name}
                        value={value}
                        setValue={(val) => handleValueChange(spice.name, val)}
                    />
                );
            case "number":
                return (
                    <NumberInputBox
                        key={spice.name}
                        value={value}
                        setValue={(val) => handleValueChange(spice.name, parseInt(val, 10))}
                    />
                );
            case "dropdown":
                return (
                    <DropdownBox
                        key={spice.name}
                        options={spice.extra}
                        selected={value}
                        onSelect={(val) => handleValueChange(spice.name, val)}
                    />
                );
            case "datepicker":
                return (
                    <DatePicker
                        key={spice.name}
                        value={value}
                        setValue={(val) => handleValueChange(spice.name, val)}
                    />
                );
            case "timestamp":
                return (
                    <TimePicker
                        key={spice.name}
                        value={value}
                        setValue={(val) => handleValueChange(spice.name, val)}
                    />
                );
            case "input":
                return (
                    <InputBox
                        key={spice.name}
                        value={value}
                        setValue={(val) => handleValueChange(spice.name, val)}
                    />
                );
            default:
                return null;
        }
    };    

    return (
        <div className="p-5 h-full flex flex-col justify-start bg-[#1d1d1d] overflow-auto font-spartan">
            <label
                className="font-bold text-2xl mb-3"
            >{feature.description}</label>
            <div className="border-b-[1px] border-chartgray-200 mb-10"></div>
            {feature.spices.map((spice, index) => (
                <div key={spice.name} className="flex flex-col mb-10">
                    <label className="text-lg">
                        {spice.title}
                    </label>
                    {renderInput(spice)}
                </div>
            ))}
            <div className="">
                <Button
                    text={"Add the new " + (feature ? "feature" : "reaction")}
                    onClick={() => {
                        const requestCheck = checkRequest(feature, request);
                        
                        if (requestCheck.length > 0) {
                            let msg = "";
                            
                            setError(true);
                            requestCheck.map((error) => {
                                msg += error + "\n";
                            });
                            setErrorMsg(msg);
                            return;
                        } else {
                            if (action) {
                                setArea((prevArea) => ({
                                    ...prevArea,
                                    actions: [...prevArea.actions, {
                                        service: service.name,
                                        name: feature.name,
                                        title: feature.description,
                                        color: service.color,
                                        spices: request
                                    }],
                                }));
                            } else {
                                setArea((prevArea) => ({
                                    ...prevArea,
                                    reactions: [...prevArea.reactions, {
                                        service: service.name,
                                        name: feature.name,
                                        title: feature.description,
                                        color: service.color,
                                        spices: request
                                    }],
                                }));
                            }

                            setError(false);
                            setErrorMsg("");
                            reset();
                        }
                    }}
                    styleClolor="bg-chartpurple-200 hover:bg-chartpurple-100 text-white"
                />
                {/* <pre className="mt-4 bg-gray-100 text-black p-4 rounded-md">
                {JSON.stringify(request, null, 2)}
               </pre>  */}
            </div>
        </div>
    );
    
}

