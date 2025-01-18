/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** FeaturesKit
*/

import Feature from "./Feature";

export default function FeaturesKit({ features, bgColor = "", setFeature, color }) {

    return (
        <div className={`grid grid-cols-1 gap-4 ${bgColor} p-5`}>

            {features && features.map((feature, index) => (
                <Feature
                    key={index}
                    feature={feature}
                    setFeature={setFeature}
                    color={color}
                />
            ))}
        </div>
    );
}