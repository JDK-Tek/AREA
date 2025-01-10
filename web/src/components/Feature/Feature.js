/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Feature
*/
export default function Feature({ id, title, color, colorHover, setFeature }) {
    return (
        <div
            className={`select-none relative text-white shadow-md p-3 cursor-pointer transition-transform duration-200`}
            onClick={() => setFeature(id)}
            style={{
                backgroundColor: color
            }}
            onMouseEnter={(e) => {
                e.currentTarget.style.backgroundColor = colorHover;
                e.currentTarget.style.transform = "scale(1.05)";
            }}
            onMouseLeave={(e) => {
                e.currentTarget.style.backgroundColor = color;
                e.currentTarget.style.transform = "scale(1)";
            }}
        >

            <h1 className="font-spartan text-xl font-bold"> {title} </h1>
        </div>
    );
}

