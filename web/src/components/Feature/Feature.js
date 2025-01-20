/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Feature
*/
export default function Feature({ feature, setFeature, color }) {
    return (
        <div
            className={`select-none relative text-white shadow-md p-3 cursor-pointer transition-transform duration-200`}
            onClick={() => setFeature(feature)}
            onKeyDown={(event) => { if (event.key === "Enter") setFeature(feature) }}
            tabIndex={0}
            style={{
                backgroundColor: color
            }}
            onMouseEnter={(e) => {
                e.currentTarget.style.backgroundColor = color;
                e.currentTarget.style.transform = "scale(1.05)";
            }}
            onMouseLeave={(e) => {
                e.currentTarget.style.backgroundColor = color;
                e.currentTarget.style.transform = "scale(1)";
            }}
        >

            <h1 className="font-spartan text-xl font-bold"> {feature.description} </h1>
        </div>
    );
}

