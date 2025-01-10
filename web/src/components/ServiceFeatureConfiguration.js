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

const simulatedConfiguration = 
        {
          "type": "reaction",
          "name": "When a new song is added to a playlist",
          "spices": [
            {
                title: "Enter the channel ID",
                name: "channel",
                type: "number",
            },
            {
                title: "Enter the message",
                name: "message",
                type: "text",
            },
            {
                title: "Select the tag",
                name: "tag",
                type: "dropdown",
                values: ["none", "everyone", "here", "Front Web"],
            },
            {
                title: "Select the date",
                name: "date",
                type: "datepicker",
            },
            {
                title: "Select the time",
                "name": "time test",
                "type": "timestamp",
            },
            {
                title: "Enter the email",
                "name": "email",
                "type": "email"
            }
          ]
        }

function checkRequest(config, request) {
    const spices = config.spices;
    const missingParameters = [];

    spices.map((spice) => {
        if (!request[spice.name]) {
            if (missingParameters.length === 0) {
                missingParameters.push("Missing parameters: ");
            }
            missingParameters.push("'" + spice.title + "', ");
        }
    });
    return missingParameters;
}

export default function ServiceFeatureConfiguration({ color, setError, setErrorMsg }) {
    const [request, setRequest] = useState({});

    const handleValueChange = (name, value) => {
        setRequest((prevRequest) => ({
            ...prevRequest,
            [name]: value,
        }));
    };

    const renderInput = (spice) => {
        const { type, name, values } = spice;
        const value = request[name] || "";

        switch (type) {
            case "text":
                return (
                    <TextBox
                        key={name}
                        value={value}
                        setValue={(val) => handleValueChange(name, val)}
                    />
                );
            case "number":
                return (
                    <NumberInputBox
                        key={name}
                        value={value}
                        setValue={(val) => handleValueChange(name, val)}
                    />
                );
            case "dropdown":
                return (
                    <DropdownBox
                        key={name}
                        options={values}
                        selected={values[0]}
                        onSelect={(val) => handleValueChange(name, val)}
                    />
                );
            case "datepicker":
                return (
                    <DatePicker
                        key={name}
                        value={value}
                        setValue={(val) => handleValueChange(name, val)}
                    />
                );
            case "timestamp":
                return (
                    <TimePicker
                        key={name}
                        value={value}
                        setValue={(val) => handleValueChange(name, val)}
                    />
                );
            case "email":
                return (
                    <InputBox
                        key={name}
                        value={value}
                        setValue={(val) => handleValueChange(name, val)}
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
                style={{ color: color.normal }}
            >{simulatedConfiguration.name}</label>
            <div className="border-b-[1px] border-chartgray-200 mb-10"></div>
            {simulatedConfiguration.spices.map((spice, index) => (
                <div key={spice.name} className="flex flex-col mb-10">
                    <label
                        className="text-lg"
                        style={{ color: color.hover }}
                    >
                        {spice.title}
                    </label>
                    {renderInput(spice)}
                </div>
            ))}
            <div className="">
                <Button
                    text="Add the new action"
                    onClick={() => {
                        const requestCheck = checkRequest(simulatedConfiguration, request);
                        
                        if (requestCheck.length > 0) {
                            let msg = "";
                            
                            setError(true);
                            requestCheck.map((error) => {
                                msg += error + "\n";
                            });
                            setErrorMsg(msg);
                            return;
                        }
                    }}
                    styleClolor="bg-chartpurple-200 hover:bg-chartpurple-100 text-white"
                />
            </div>
            <pre className="mt-4 bg-gray-100 text-black p-4 rounded-md">
                {JSON.stringify(request, null, 2)}
            </pre>
        </div>
    );
    
}

