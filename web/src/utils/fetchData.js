/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** fetchData
*/

export default async function fetchData(url, request = {method: "GET"}) {

    try {
        const res = await fetch(url, request);
        if (!res.ok) {
            console.error(`Response status: ${res}`);
            throw new Error(`Response status: ${res}`);
        }
    
        const json = await res.json();
        return { success: true, data: json };
    } catch (err) {
        console.error(`Error while fetching data: ${err}`);
        return { success: false, error: err };
    }
}
