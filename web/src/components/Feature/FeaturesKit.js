/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** FeaturesKit
*/

import Feature from "./Feature";

export default function FeaturesKit({ features }) {
    return (
        <div className="grid grid-cols-1 gap-4">
            {features.features.map((feature, index) => (
                <Feature key={index} title={feature.title} color={features.color} colorHover={features.colorHover} />
            ))}
        </div>
    );
}